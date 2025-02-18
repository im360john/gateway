package postgres

import (
	"context"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func init() {
	connectors.Register[Config](func(cfg Config) (connectors.Connector, error) {
		return &Connector{
			config: cfg,
		}, nil
	})
}

type Connector struct {
	config Config
}

func (c Connector) Ping(ctx context.Context) error {
	cfg, err := c.config.MakeConfig()
	if err != nil {
		return errors.Errorf("unable to prepare pg config: %w", err)
	}
	db := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")
	rows, err := db.QueryContext(ctx, "select 1+1")
	if err != nil {
		return errors.Errorf("unable to ping db: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var res int
		if err := rows.Scan(&res); err != nil {
			return errors.Errorf("unable to scan ping result: %w", err)
		}
	}
	if rows.Err() != nil {
		return errors.Errorf("rows fetcher failed: %w", err)
	}
	return nil
}

func (c Connector) Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error) {
	cfg, err := c.config.MakeConfig()
	if err != nil {
		return nil, errors.Errorf("unable to prepare pg config: %w", err)
	}
	db := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")
	rows, err := db.NamedQueryContext(ctx, endpoint.Query, params)
	if err != nil {
		return nil, errors.Errorf("unable to ping db: %w", err)
	}
	defer rows.Close()
	res := make([]map[string]any, 0)
	for rows.Next() {
		row := map[string]any{}
		if err := rows.MapScan(row); err != nil {
			return nil, errors.Errorf("unable to scan row: %w", err)
		}
		res = append(res, row)
	}
	return res, nil
}
