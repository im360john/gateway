package clickhouse

import (
	"context"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/centralmind/gateway/castx"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/xerrors"
)

func init() {
	connectors.Register(func(cfg Config) (connectors.Connector, error) {
		dsn := cfg.MakeDSN()
		db, err := sqlx.Open("clickhouse", dsn)
		if err != nil {
			return nil, xerrors.Errorf("unable to open ClickHouse db: %w", err)
		}
		return &Connector{
			config: cfg,
			db:     db,
		}, nil
	})
}

type Config struct {
	Host     string   // Single host
	Hosts    []string // Multiple hosts
	Database string
	User     string
	Password string
	Port     int
	Secure   bool
}

func (c Config) MakeDSN() string {
	protocol := "http"
	if c.Secure {
		protocol = "https"
	}

	host := c.Host
	// If no single host is specified but we have hosts array, use the first one
	if host == "" && len(c.Hosts) > 0 {
		host = c.Hosts[0]
	}

	// Format as protocol://user:password@host:port/database
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", protocol, c.User, c.Password, host, c.Port, c.Database)
}

func (c Config) Type() string {
	return "clickhouse"
}

func (c Config) Doc() string {
	return `ClickHouse connector allows querying ClickHouse databases.

Config example:
    host: localhost      # Single host address
    hosts:              # Or multiple hosts for cluster setup
      - host1.example.com
      - host2.example.com
    database: mydb
    user: default
    password: secret
    port: 8123
    secure: false       # Use HTTPS instead of HTTP`
}

type Connector struct {
	config Config
	db     *sqlx.DB
}

func (c Connector) Sample(ctx context.Context, table model.Table) ([]map[string]any, error) {
	rows, err := c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s LIMIT 5", table.Name))
	if err != nil {
		return nil, xerrors.Errorf("unable to query db: %w", err)
	}
	defer rows.Close()

	res := make([]map[string]any, 0, 5)
	for rows.Next() {
		row := make(map[string]any)
		if err := rows.MapScan(row); err != nil {
			return nil, xerrors.Errorf("unable to scan row: %w", err)
		}
		res = append(res, row)
	}
	return res, nil
}

func (c Connector) Discovery(ctx context.Context) ([]model.Table, error) {
	rows, err := c.db.Query("SELECT name FROM system.tables WHERE database = ?", c.config.Database)
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
	return c.db.PingContext(ctx)
}

func (c Connector) Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error) {
	processed, err := castx.ParamsE(endpoint, params)
	if err != nil {
		return nil, xerrors.Errorf("unable to process params: %w", err)
	}
	rows, err := c.db.NamedQueryContext(ctx, endpoint.Query, processed)
	if err != nil {
		return nil, xerrors.Errorf("unable to query db: %w", err)
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
		"SELECT name, type FROM system.columns WHERE table = ? AND database = ?",
		tableName, c.config.Database,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []model.ColumnSchema
	for rows.Next() {
		var name, dataType string
		if err := rows.Scan(&name, &dataType); err != nil {
			return nil, err
		}
		columns = append(columns, model.ColumnSchema{
			Name: name,
			Type: dataType,
		})
	}
	return columns, nil
}
