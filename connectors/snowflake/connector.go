package snowflake

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/centralmind/gateway/connectors"

	"github.com/centralmind/gateway/castx"
	"github.com/centralmind/gateway/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/snowflakedb/gosnowflake"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

//go:embed readme.md
var docString string

func init() {
	connectors.Register(func(cfg Config) (connectors.Connector, error) {
		dsn, err := cfg.MakeDSN()
		if err != nil {
			return nil, xerrors.Errorf("unable to prepare Snowflake config: %w", err)
		}
		db, err := sqlx.Open("snowflake", dsn)
		if err != nil {
			return nil, xerrors.Errorf("unable to open Snowflake db: %w", err)
		}
		return &Connector{
			config: cfg,
			db:     db,
			base:   &connectors.BaseConnector{DB: db},
		}, nil
	})
}

type Config struct {
	Account    string
	Database   string
	User       string
	Password   string
	Warehouse  string
	Schema     string
	Role       string
	ConnString string `yaml:"conn_string"`
	IsReadonly bool   `yaml:"is_readonly"`
}

func (c Config) Readonly() bool {
	return c.IsReadonly
}

// UnmarshalYAML implements the yaml.Unmarshaler interface to allow for both
// direct connection string or full configuration objects in YAML
func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal as a string (connection string)
	var connString string
	if err := value.Decode(&connString); err == nil && len(connString) > 0 {
		c.ConnString = connString
		return nil
	}

	// If that didn't work, try to unmarshal as a full config object
	type configAlias Config // Use alias to avoid infinite recursion
	var alias configAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}

	*c = Config(alias)
	return nil
}

func (c Config) ExtraPrompt() []string {
	return []string{}
}

func (c Config) MakeDSN() (string, error) {
	// If connection string is provided, use it directly
	if c.ConnString != "" {
		return c.ConnString, nil
	}

	// Otherwise, build the DSN from individual fields
	dsn := fmt.Sprintf("%s:%s@%s/%s/%s?warehouse=%s&role=%s", c.User, c.Password, c.Account, c.Database, c.Schema, c.Warehouse, c.Role)

	return dsn, nil
}

func (c Config) Type() string {
	return "snowflake"
}

func (c Config) Doc() string {
	return docString
}

type Connector struct {
	config Config
	db     *sqlx.DB
	base   *connectors.BaseConnector
}

func (c Connector) Config() connectors.Config {
	return c.config
}

func (c Connector) Sample(ctx context.Context, table model.Table) ([]map[string]any, error) {
	rows, err := c.db.NamedQuery(fmt.Sprintf("SELECT * FROM %s LIMIT 5", table.Name), map[string]any{})
	if err != nil {
		return nil, xerrors.Errorf("unable to query db: %w", err)
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

func (c Connector) Discovery(ctx context.Context, tablesList []string) ([]model.Table, error) {
	// Create a map for quick lookups if tablesList is provided
	tableSet := make(map[string]bool)
	if len(tablesList) > 0 {
		for _, table := range tablesList {
			tableSet[table] = true
		}
	}

	// Create base query
	queryBase := fmt.Sprintf("SHOW TABLES IN SCHEMA %s.%s", c.config.Database, c.config.Schema)

	// If specific tables are requested, filter with LIKE conditions
	// Note: Snowflake doesn't support WHERE IN for SHOW TABLES, so we need to use individual queries
	var allTables []model.Table

	if len(tablesList) > 0 {
		// For each requested table, query individually
		for _, tableName := range tablesList {
			// Use LIKE to match the exact table name
			query := queryBase + fmt.Sprintf(" LIKE '%s'", tableName)
			tables, err := c.executeTableQuery(ctx, query)
			if err != nil {
				return nil, err
			}
			allTables = append(allTables, tables...)
		}
	} else {
		// If no specific tables are requested, get all tables
		tables, err := c.executeTableQuery(ctx, queryBase)
		if err != nil {
			return nil, err
		}
		allTables = tables
	}

	return allTables, nil
}

// Helper function to execute table queries and process results
func (c Connector) executeTableQuery(ctx context.Context, query string) ([]model.Table, error) {
	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []model.Table
	for rows.Next() {
		var tableName string
		var createdOn, kind, databaseName, schemaName, clusterBy, owner, comment, changeTracking, automaticClustering, searchOptimization string

		if err := rows.Scan(&createdOn, &tableName, &kind, &databaseName, &schemaName, &clusterBy, &owner, &comment, &changeTracking, &automaticClustering, &searchOptimization); err != nil {
			return nil, err
		}
		columns, err := c.LoadsColumns(ctx, tableName)
		if err != nil {
			return nil, err
		}

		// Get the total row count for this table
		var rowCount int
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\".\"%s\".\"%s\"", c.config.Database, c.config.Schema, tableName)
		err = c.db.Get(&rowCount, countQuery)
		if err != nil {
			return nil, xerrors.Errorf("unable to get row count for table %s: %w", tableName, err)
		}

		table := model.Table{
			Name:     tableName,
			Columns:  columns,
			RowCount: rowCount,
		}
		tables = append(tables, table)
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

	rows, err := c.db.NamedQuery(endpoint.Query, processed)
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
		`SELECT 
			c.COLUMN_NAME,
			c.DATA_TYPE,
			CASE WHEN k.COLUMN_NAME IS NOT NULL THEN true ELSE false END as is_primary_key
		FROM information_schema.columns c
		LEFT JOIN information_schema.key_column_usage k 
			ON c.table_catalog = k.table_catalog 
			AND c.table_schema = k.table_schema
			AND c.table_name = k.table_name 
			AND c.column_name = k.column_name 
			AND k.constraint_name LIKE 'SYS_CONSTRAINT_%'
		WHERE c.table_name = ?
		AND c.table_schema = ?
		AND c.table_catalog = ?`,
		tableName, c.config.Schema, c.config.Database,
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to query columns: %w", err)
	}
	defer rows.Close()

	var columns []model.ColumnSchema
	for rows.Next() {
		var name, dataType string
		var isPrimaryKey bool
		if err := rows.Scan(&name, &dataType, &isPrimaryKey); err != nil {
			return nil, err
		}
		columns = append(columns, model.ColumnSchema{
			Name:       name,
			Type:       c.GuessColumnType(dataType),
			PrimaryKey: isPrimaryKey,
		})
	}
	return columns, nil
}

// GuessColumnType implements TypeGuesser interface for Snowflake
func (c *Connector) GuessColumnType(sqlType string) model.ColumnType {
	upperType := strings.ToUpper(sqlType)

	// Array type
	if strings.Contains(upperType, "ARRAY") {
		return model.TypeArray
	}

	// Object types
	switch upperType {
	case "OBJECT", "VARIANT":
		return model.TypeObject
	}

	// String types
	switch upperType {
	case "STRING", "TEXT", "VARCHAR", "CHAR", "BINARY", "VARBINARY":
		return model.TypeString
	}

	// Numeric types
	switch upperType {
	case "NUMBER", "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE":
		return model.TypeNumber
	}

	// Integer types
	switch upperType {
	case "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT":
		return model.TypeInteger
	}

	// Boolean type
	if upperType == "BOOLEAN" {
		return model.TypeBoolean
	}

	// Date/Time types
	switch upperType {
	case "DATE", "TIME", "TIMESTAMP", "TIMESTAMP_LTZ", "TIMESTAMP_NTZ", "TIMESTAMP_TZ":
		return model.TypeDatetime
	}

	// Default to string for unknown types
	return model.TypeString
}

// InferResultColumns returns column information for the given query
func (c *Connector) InferResultColumns(ctx context.Context, query string) ([]model.ColumnSchema, error) {
	return c.base.InferResultColumns(ctx, query, c)
}

// InferQuery implements the Connector interface
func (c *Connector) InferQuery(ctx context.Context, query string) ([]model.ColumnSchema, error) {
	return c.base.InferResultColumns(ctx, query, c)
}
