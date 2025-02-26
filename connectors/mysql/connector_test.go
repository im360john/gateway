package mysql

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestConnector(t *testing.T) {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUser),
		mysql.WithPassword(dbPassword),
		mysql.WithScripts(filepath.Join("testdata", "test_data.sql")),
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testcontainers.TerminateContainer(mysqlContainer))
	}()

	host, err := mysqlContainer.Host(ctx)
	require.NoError(t, err)
	port, err := mysqlContainer.MappedPort(ctx, nat.Port("3306/tcp"))
	require.NoError(t, err)

	logrus.Infof("mysql running on %s:%d", host, port.Int())

	cfg := Config{
		Host:     host,
		Database: dbName,
		User:     dbUser,
		Password: dbPassword,
		Port:     port.Int(),
	}

	connector, err := connectors.New("mysql", cfg)
	assert.NoError(t, err)
	assert.NotNil(t, connector)

	t.Run("Ping Database", func(t *testing.T) {
		err := connector.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Discovery Tables", func(t *testing.T) {
		tables, err := connector.Discovery(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, tables)
	})

	t.Run("Read Endpoint", func(t *testing.T) {
		endpoint := model.Endpoint{
			Query:  "SELECT COUNT(*) AS total_count FROM gachi_teams",
			Params: []model.EndpointParams{},
		}
		params := map[string]any{}
		rows, err := connector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)
	})

	t.Run("Query Endpoint", func(t *testing.T) {
		endpoint := model.Endpoint{
			Query: `SELECT name, strength_level, special_move 
					FROM gachi_personas 
					WHERE strength_level > :min_strength 
					ORDER BY strength_level DESC`,
			Params: []model.EndpointParams{
				{
					Name:     "min_strength",
					Type:     "int",
					Required: true,
				},
			},
		}
		params := map[string]any{
			"min_strength": 90,
		}
		rows, err := connector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)

		// Verify results
		expected := []map[string]any{
			{"name": "Billy Herrington", "strength_level": int64(100), "special_move": "Anvil Drop"},
			{"name": "Van Darkholme", "strength_level": int64(95), "special_move": "Whip of Submission"},
		}

		assert.Equal(t, len(expected), len(rows))
		for i, exp := range expected {
			assert.Equal(t, exp["name"], rows[i]["name"])
			assert.Equal(t, exp["strength_level"], rows[i]["strength_level"])
			assert.Equal(t, exp["special_move"], rows[i]["special_move"])
		}
	})
}
