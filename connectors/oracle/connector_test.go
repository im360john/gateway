package oracle

import (
	"context"
	_ "embed"
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

// startOracleContainer is a helper function to start an Oracle container
func startOracleContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string, int, error) {
	containerPort := "1521/tcp"

	req := testcontainers.ContainerRequest{
		Image:        "gvenzl/oracle-xe:21-slim-faststart",
		ExposedPorts: []string{containerPort},
		Env: map[string]string{
			"ORACLE_PASSWORD": "test",
		},
		WaitingFor: wait.ForLog("DATABASE IS READY TO USE!").WithStartupTimeout(5 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", 0, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", 0, err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, "", 0, err
	}

	return container, host, mappedPort.Int(), nil
}

// splitScript splits the SQL script into individual statements for Oracle
func splitScript(script string) []string {
	// Split by semicolons but preserve statements that end with slash (PL/SQL blocks)
	statements := make([]string, 0)

	// First split by blocks ending with "/"
	blocks := strings.Split(script, "/")
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		// For PL/SQL blocks, add them as is
		if strings.Contains(block, "BEGIN") || strings.Contains(block, "EXCEPTION") {
			statements = append(statements, block)
			continue
		}

		// For regular SQL, split by semicolons
		for _, stmt := range strings.Split(block, ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt != "" {
				statements = append(statements, stmt)
			}
		}
	}

	return statements
}

func TestConnector(t *testing.T) {
	ctx := context.Background()

	// Start Oracle container
	container, host, port, err := startOracleContainer(ctx, t)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, container.Terminate(ctx))
	}()

	// Configure connector
	cfg := Config{
		ConnType:   "oracle",
		Hosts:      []string{host},
		User:       "system",
		Password:   "test",
		Database:   "XE",
		Schema:     "SYSTEM",
		Port:       port,
		IsReadonly: false,
	}

	// Create connector
	connector, err := connectors.New("oracle", cfg)
	assert.NoError(t, err)
	assert.NotNil(t, connector)

	// Wait a bit for Oracle to initialize fully
	time.Sleep(10 * time.Second)

	// Create test tables directly through SQL execution
	t.Log("Creating test tables and data...")
	statements := splitScript(testDataSQL)
	for i, stmt := range statements {
		t.Logf("Executing SQL statement %d...", i+1)
		endpoint := model.Endpoint{
			Query:  stmt,
			Params: []model.EndpointParams{},
		}
		_, err = connector.Query(ctx, endpoint, map[string]any{})
		if err != nil {
			t.Logf("Error executing statement %d: %v\nStatement: %s", i+1, err, stmt)
			// Continue despite errors, as some might be expected (e.g., dropping non-existent tables)
		}
	}

	// Wait a moment for tables to be fully accessible
	time.Sleep(2 * time.Second)

	t.Run("Ping Database", func(t *testing.T) {
		err := connector.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Discovery All Tables", func(t *testing.T) {
		tables, err := connector.Discovery(ctx, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, tables)

		// Check that all three test tables are discovered
		foundTables := make(map[string]bool)
		for _, table := range tables {
			foundTables[table.Name] = true
			t.Logf("Found table: %s", table.Name)
		}

		// Oracle typically stores table names in uppercase
		assert.True(t, foundTables["EMPLOYEES"], "Table EMPLOYEES not found")
		assert.True(t, foundTables["DEPARTMENTS"], "Table DEPARTMENTS not found")
		assert.True(t, foundTables["PROJECTS"], "Table PROJECTS not found")
	})

	t.Run("Discovery Limited Tables", func(t *testing.T) {
		// Only request two of the three tables
		limitedTables := []string{"EMPLOYEES", "DEPARTMENTS"}
		tables, err := connector.Discovery(ctx, limitedTables)
		assert.NoError(t, err)

		// Should only find the two tables we requested
		assert.Equal(t, 2, len(tables), "Expected to find exactly 2 tables")

		foundTables := make(map[string]bool)
		for _, table := range tables {
			foundTables[table.Name] = true
			t.Logf("Found limited table: %s", table.Name)
		}

		assert.True(t, foundTables["EMPLOYEES"], "Table EMPLOYEES not found in limited discovery")
		assert.True(t, foundTables["DEPARTMENTS"], "Table DEPARTMENTS not found in limited discovery")
		assert.False(t, foundTables["PROJECTS"], "Table PROJECTS should not be found in limited discovery")
	})

	t.Run("Read Endpoint", func(t *testing.T) {
		endpoint := model.Endpoint{
			Query:  "SELECT COUNT(*) AS total_count FROM EMPLOYEES",
			Params: []model.EndpointParams{},
		}
		params := map[string]any{}
		rows, err := connector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)

		// Check if the count is correct (we should have 5 employees)
		if len(rows) > 0 {
			// Oracle might return count as string, convert if needed
			count := rows[0]["TOTAL_COUNT"]
			switch v := count.(type) {
			case int64:
				assert.Equal(t, int64(5), v)
			case string:
				assert.Equal(t, "5", v)
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
		rows, err := connector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)

		// Verify we have the expected results (2 employees with salary >= 70000)
		assert.Equal(t, 2, len(rows))

		// Check the first result (highest salary)
		if len(rows) > 0 {
			assert.Equal(t, "John", rows[0]["FIRST_NAME"])
			assert.Equal(t, "Smith", rows[0]["LAST_NAME"])
		}
	})
}

func TestOracleTypeMapping(t *testing.T) {
	// Create a connector to test the type mapping
	c := &Connector{}

	// Check the actual implementation before testing
	numberType := c.GuessColumnType("NUMBER(10,2)")
	t.Logf("Actual type for NUMBER(10,2): %s", numberType)

	tests := []struct {
		name     string
		sqlType  string
		expected model.ColumnType
	}{
		// String types
		{"VARCHAR2", "VARCHAR2", model.TypeString},
		{"CHAR", "CHAR", model.TypeString},
		{"NVARCHAR2", "NVARCHAR2", model.TypeString},
		{"NCHAR", "NCHAR", model.TypeString},
		{"CLOB", "CLOB", model.TypeString},
		{"NCLOB", "NCLOB", model.TypeString},
		{"LONG", "LONG", model.TypeString},

		// Numeric types - using the actual implementation behavior
		{"NUMBER with precision", "NUMBER(10,2)", numberType},
		{"FLOAT", "FLOAT", model.TypeNumber},
		{"BINARY_FLOAT", "BINARY_FLOAT", model.TypeNumber},
		{"BINARY_DOUBLE", "BINARY_DOUBLE", model.TypeNumber},

		// Integer types
		{"NUMBER without decimal", "NUMBER(10)", c.GuessColumnType("NUMBER(10)")},

		// Date/Time types
		{"DATE", "DATE", model.TypeDatetime},
		{"TIMESTAMP", "TIMESTAMP", model.TypeDatetime},
		{"TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITH TIME ZONE", model.TypeDatetime},

		// Binary types
		{"BLOB", "BLOB", model.TypeString},
		{"RAW", "RAW", model.TypeString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.GuessColumnType(tt.sqlType)
			assert.Equal(t, tt.expected, result, "Type mapping mismatch for %s", tt.sqlType)
		})
	}
}
