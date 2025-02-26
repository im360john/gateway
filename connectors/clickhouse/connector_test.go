package clickhouse

import (
	"context"
	_ "embed"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/clickhouse"
	"path/filepath"
	"testing"
)

func TestConnector(t *testing.T) {
	ctx := context.Background()
	user := "clickhouse"
	password := "password"
	dbname := "testdb"
	clickHouseContainer, err := clickhouse.Run(ctx,
		"clickhouse/clickhouse-server:23.3.8.21-alpine",
		clickhouse.WithUsername(user),
		clickhouse.WithPassword(password),
		clickhouse.WithDatabase(dbname),
		clickhouse.WithInitScripts(filepath.Join("testdata", "test_data.sql")),
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testcontainers.TerminateContainer(clickHouseContainer))
	}()

	connStr, err := clickHouseContainer.ConnectionString(ctx)
	require.NoError(t, err)
	logrus.Infof("conn string: %s", connStr)
	host, err := clickHouseContainer.Host(ctx)
	require.NoError(t, err)
	port, err := clickHouseContainer.MappedPort(ctx, nat.Port("8123/tcp"))
	require.NoError(t, err)
	cfg := Config{
		Host:     host,
		Database: dbname,
		User:     user,
		Password: password,
		Port:     port.Int(),
		Secure:   false,
	}
	connector, err := connectors.New("clickhouse", cfg)
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
					WHERE strength_level >= :min_strength 
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
			"min_strength": 95,
		}
		rows, err := connector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)

		// Verify results
		expected := []map[string]any{
			{"name": "Billy Herrington", "strength_level": int32(100), "special_move": "Anvil Drop"},
			{"name": "Van Darkholme", "strength_level": int32(95), "special_move": "Whip of Submission"},
		}

		assert.Equal(t, len(expected), len(rows))
		for i, exp := range expected {
			assert.Equal(t, exp["name"], rows[i]["name"])
			assert.Equal(t, exp["strength_level"], rows[i]["strength_level"])
			assert.Equal(t, exp["special_move"], rows[i]["special_move"])
		}
	})
}
