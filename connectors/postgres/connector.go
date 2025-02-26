package postgres

import (
	"context"
	"fmt"
	"github.com/centralmind/gateway/castx"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/xerrors"
)

func init() {
	connectors.Register[Config](func(cfg Config) (connectors.Connector, error) {
		c, err := cfg.MakeConfig()
		if err != nil {
			return nil, xerrors.Errorf("unable to prepare pg config: %w", err)
		}
		db := sqlx.NewDb(stdlib.OpenDB(*c), "pgx")
		return &Connector{
			config: cfg,
			db:     db,
		}, nil
	})
}

type Connector struct {
	config Config
	db     *sqlx.DB
}

func (c Connector) Sample(ctx context.Context, table model.Table) ([]map[string]any, error) {
	rows, err := c.db.NamedQueryContext(ctx, fmt.Sprintf("select * from %s limit 5", table.Name), map[string]any{})
	if err != nil {
		return nil, xerrors.Errorf("unable to ping db: %w", err)
	}
	defer rows.Close()
	res := make([]map[string]any, 0, 5)
	for rows.Next() {
		row := map[string]any{}
		if err := rows.MapScan(row); err != nil {
			return nil, xerrors.Errorf("unable to scan row: %w", err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (c Connector) Discovery(ctx context.Context) ([]model.Table, error) {
	cfg, err := c.config.MakeConfig()
	if err != nil {
		return nil, xerrors.Errorf("unable to prepare pg config: %w", err)
	}
	db := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []model.Table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		columns, err := c.LoadsColumns(ctx, tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, model.Table{Name: tableName, Columns: columns})
	}
	return tables, nil
}

func (c Connector) Ping(ctx context.Context) error {
	rows, err := c.db.QueryContext(ctx, "select 1+1")
	if err != nil {
		return xerrors.Errorf("unable to ping db: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var res int
		if err := rows.Scan(&res); err != nil {
			return xerrors.Errorf("unable to scan ping result: %w", err)
		}
	}
	if rows.Err() != nil {
		return xerrors.Errorf("rows fetcher failed: %w", err)
	}
	return nil
}

func (c Connector) Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error) {
	processed, err := castx.ParamsE(endpoint, params)
	if err != nil {
		return nil, xerrors.Errorf("unable to process params: %w", err)
	}

	rows, err := c.db.NamedQueryContext(ctx, endpoint.Query, processed)
	if err != nil {
		return nil, xerrors.Errorf("unable to ping db: %w", err)
	}
	defer rows.Close()
	res := make([]map[string]any, 0)
	for rows.Next() {
		row := map[string]any{}
		if err := rows.MapScan(row); err != nil {
			return nil, xerrors.Errorf("unable to scan row: %w", err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (c Connector) LoadsColumns(ctx context.Context, tableName string) ([]model.ColumnSchema, error) {
	rows, err := c.db.QueryContext(
		ctx,
		`SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = $1`,
		tableName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []model.ColumnSchema
	for rows.Next() {
		var name, dataType, isNullable string
		if err := rows.Scan(&name, &dataType, &isNullable); err != nil {
			return nil, err
		}
		columns = append(columns, model.ColumnSchema{
			Name: name,
			Type: dataType,
		})
	}
	return columns, nil
}
