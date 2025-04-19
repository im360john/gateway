package mssql

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

//go:embed testdata/test_data.sql
var testDataSQL string

// startMSSQLContainer is a helper function to start a MSSQL container
func startMSSQLContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string, int, string, error) {
	containerPort := "1433/tcp"
	// Create a password that definitely meets SQL Server 2022 complexity requirements
	password := "StrongPwd1!" // Has uppercase, lowercase, number, and symbol, length > 8

	req := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/mssql/server:2022-latest", // Use SQL Server 2022
		ExposedPorts: []string{containerPort},
		Env: map[string]string{
			"ACCEPT_EULA":       "Y",
			"MSSQL_SA_PASSWORD": password, // SQL Server 2022 environment variable
			"MSSQL_PID":         "Developer",
		},
		WaitingFor: wait.ForLog("SQL Server is now ready for client connections").WithStartupTimeout(5 * time.Minute),
	}

	// Print the environment variables for debugging
	t.Logf("SQL Server environment: ACCEPT_EULA=Y, MSSQL_SA_PASSWORD=%s, MSSQL_PID=Developer", password)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", 0, "", err
	}

	// Wait for SQL Server to be fully initialized
	time.Sleep(3 * time.Second)

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", 0, "", err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, "", 0, "", err
	}

	return container, host, mappedPort.Int(), password, nil
}

// splitScript splits the SQL script into individual statements for MSSQL
func splitScript(script string) []string {
	statements := make([]string, 0)
	for _, stmt := range strings.Split(script, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	return statements
}

func TestConnector(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MSSQL integration test in short mode")
	}

	ctx := context.Background()

	// Start MSSQL container
	t.Log("Starting MSSQL container...")
	container, host, port, password, err := startMSSQLContainer(ctx, t)
	require.NoError(t, err)
	defer func() {
		t.Log("Terminating MSSQL container...")
		require.NoError(t, container.Terminate(ctx))
	}()

	t.Logf("MSSQL container started at %s:%d", host, port)

	// Configure connector
	cfg := Config{
		Hosts:    []string{host},
		User:     "sa",
		Password: password,
		Database: "master", // Use master initially
		Port:     port,
		Schema:   "dbo",
	}
	t.Logf("Connection config: %+v", cfg)

	// Log the connection string for debugging
	connStr := cfg.ConnectionString()
	t.Logf("Connection string: %s", connStr)

	// Create connector
	t.Log("Creating connector to master database...")
	connector, err := connectors.New("mssql", cfg)
	if err != nil {
		t.Fatalf("Failed to create connector: %v", err)
	}
	require.NotNil(t, connector)

	// Test ping to master database
	t.Log("Pinging master database...")
	err = connector.Ping(ctx)
	require.NoError(t, err, "Failed to ping master database")

	// Wait for SQL Server to be fully ready
	t.Log("Waiting for SQL Server to be fully ready...")
	time.Sleep(5 * time.Second)

	// Create test database
	dbName := "TestDB"
	t.Logf("Creating test database: %s", dbName)
	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
	createDBEndpoint := model.Endpoint{
		Query:  createDBQuery,
		Params: []model.EndpointParams{},
	}
	_, err = connector.Query(ctx, createDBEndpoint, map[string]any{})
	require.NoError(t, err, "Failed to create test database")

	// Wait a moment for the database to be fully created
	time.Sleep(2 * time.Second)

	// Now switch to the test database
	cfg.Database = dbName
	t.Logf("Switching to test database: %s", dbName)
	testDBConnector, err := connectors.New("mssql", cfg)
	require.NoError(t, err)
	require.NotNil(t, testDBConnector)

	// Create test tables and insert data as a single batch for better reliability
	t.Log("Creating test tables and data...")

	// First create tables one by one
	statements := splitScript(testDataSQL)

	// Process DROP TABLE statements first
	for i, stmt := range statements[:3] { // The first 3 statements in our SQL are DROP TABLE statements
		t.Logf("Executing DROP statement %d: %s", i+1, stmt)
		endpoint := model.Endpoint{
			Query:  stmt,
			Params: []model.EndpointParams{},
		}
		_, err = testDBConnector.Query(ctx, endpoint, map[string]any{})
		// Ignore errors from DROP statements as tables may not exist yet
		if err != nil {
			t.Logf("Expected error dropping table: %v", err)
		}
	}

	// Then execute CREATE TABLE statements
	for i, stmt := range statements[3:6] { // Statements 3-5 are CREATE TABLE statements
		t.Logf("Executing CREATE TABLE statement %d: %s", i+3+1, stmt)
		endpoint := model.Endpoint{
			Query:  stmt,
			Params: []model.EndpointParams{},
		}
		_, err = testDBConnector.Query(ctx, endpoint, map[string]any{})
		require.NoError(t, err, "Failed to create table with statement: %s", stmt)
	}

	// Finally execute INSERT statements
	for i, stmt := range statements[6:] { // Remaining statements are INSERTs
		t.Logf("Executing INSERT statement %d: %s", i+6+1, stmt)
		endpoint := model.Endpoint{
			Query:  stmt,
			Params: []model.EndpointParams{},
		}
		_, err = testDBConnector.Query(ctx, endpoint, map[string]any{})
		require.NoError(t, err, "Failed to insert data with statement: %s", stmt)
	}

	t.Run("Ping Database", func(t *testing.T) {
		err := testDBConnector.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Discovery All Tables", func(t *testing.T) {
		tables, err := testDBConnector.Discovery(ctx, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, tables)

		// Check that all three test tables are discovered
		foundTables := make(map[string]bool)
		for _, table := range tables {
			foundTables[table.Name] = true
			t.Logf("Found table: %s", table.Name)
		}

		assert.True(t, foundTables["employees"], "Table employees not found")
		assert.True(t, foundTables["departments"], "Table departments not found")
		assert.True(t, foundTables["projects"], "Table projects not found")
	})

	t.Run("Discovery Limited Tables", func(t *testing.T) {
		// Only request two of the three tables
		limitedTables := []string{"employees", "departments"}
		tables, err := testDBConnector.Discovery(ctx, limitedTables)
		assert.NoError(t, err)

		// Should only find the two tables we requested
		assert.Equal(t, 2, len(tables), "Expected to find exactly 2 tables")

		foundTables := make(map[string]bool)
		for _, table := range tables {
			foundTables[table.Name] = true
			t.Logf("Found limited table: %s", table.Name)
		}

		assert.True(t, foundTables["employees"], "Table employees not found in limited discovery")
		assert.True(t, foundTables["departments"], "Table departments not found in limited discovery")
		assert.False(t, foundTables["projects"], "Table projects should not be found in limited discovery")
	})

	t.Run("Read Endpoint", func(t *testing.T) {
		endpoint := model.Endpoint{
			Query:  "SELECT COUNT(*) AS total_count FROM employees",
			Params: []model.EndpointParams{},
		}
		params := map[string]any{}
		rows, err := testDBConnector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)

		// Check if the count is correct (we should have 5 employees)
		if len(rows) > 0 {
			// Check the type and value
			count := rows[0]["total_count"]
			switch v := count.(type) {
			case int64:
				assert.Equal(t, int64(5), v)
			case int32:
				assert.Equal(t, int32(5), v)
			case int:
				assert.Equal(t, 5, v)
			default:
				t.Logf("Unexpected type for count: %T", count)
				assert.Equal(t, 5, count) // Will fail with details about the actual type
			}
		}
	})

	t.Run("Query Endpoint With Params", func(t *testing.T) {
		endpoint := model.Endpoint{
			Query: `SELECT first_name, last_name, salary 
					FROM employees 
					WHERE salary >= :min_salary 
					ORDER BY salary DESC`,
			Params: []model.EndpointParams{
				{
					Name:     "min_salary",
					Type:     "number",
					Required: true,
				},
			},
		}
		params := map[string]any{
			"min_salary": 70000,
		}
		rows, err := testDBConnector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)

		// Verify we have the expected results (2 employees with salary >= 70000)
		assert.Equal(t, 2, len(rows))

		// Check the first result (highest salary)
		if len(rows) > 0 {
			assert.Equal(t, "John", rows[0]["first_name"])
			assert.Equal(t, "Smith", rows[0]["last_name"])
		}
	})
}

func TestMSSQLTypeMapping(t *testing.T) {
	c := &Connector{}

	tests := []struct {
		name     string
		sqlType  string
		expected model.ColumnType
	}{
		// String types
		{"VARCHAR", "VARCHAR", model.TypeString},
		{"CHAR", "CHAR", model.TypeString},
		{"TEXT", "TEXT", model.TypeString},
		{"NVARCHAR", "NVARCHAR", model.TypeString},
		{"NCHAR", "NCHAR", model.TypeString},
		{"NTEXT", "NTEXT", model.TypeString},
		{"XML", "XML", model.TypeString},
		{"UNIQUEIDENTIFIER", "UNIQUEIDENTIFIER", model.TypeString},

		// Numeric types
		{"DECIMAL", "DECIMAL", model.TypeNumber},
		{"NUMERIC", "NUMERIC", model.TypeNumber},
		{"FLOAT", "FLOAT", model.TypeNumber},
		{"REAL", "REAL", model.TypeNumber},
		{"MONEY", "MONEY", model.TypeNumber},
		{"SMALLMONEY", "SMALLMONEY", model.TypeNumber},

		// Integer types
		{"INT", "INT", model.TypeInteger},
		{"BIGINT", "BIGINT", model.TypeInteger},
		{"SMALLINT", "SMALLINT", model.TypeInteger},
		{"TINYINT", "TINYINT", model.TypeInteger},

		// Boolean type
		{"BIT", "BIT", model.TypeBoolean},

		// Date/Time types
		{"DATE", "DATE", model.TypeDatetime},
		{"TIME", "TIME", model.TypeDatetime},
		{"DATETIME", "DATETIME", model.TypeDatetime},
		{"DATETIME2", "DATETIME2", model.TypeDatetime},
		{"DATETIMEOFFSET", "DATETIMEOFFSET", model.TypeDatetime},
		{"SMALLDATETIME", "SMALLDATETIME", model.TypeDatetime},

		// Binary types
		{"BINARY", "BINARY", model.TypeString},
		{"VARBINARY", "VARBINARY", model.TypeString},
		{"IMAGE", "IMAGE", model.TypeString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.GuessColumnType(tt.sqlType)
			assert.Equal(t, tt.expected, result, "Type mapping mismatch for %s", tt.sqlType)
		})
	}
}
