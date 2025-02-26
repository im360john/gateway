package postgres

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestConnector(t *testing.T) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("testdata", "test_data.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testcontainers.TerminateContainer(postgresContainer))
	}()

	connStr, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)
	logrus.Infof("conn string: %s", connStr)
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, nat.Port("5432/tcp"))
	require.NoError(t, err)
	cfg := Config{
		Hosts:    []string{host},
		Database: dbName,
		User:     dbUser,
		Password: dbPassword,
		Port:     port.Int(),
	}
	connector, err := connectors.New("postgres", cfg)
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
			Query: `SELECT name, strength_level, special_move, team_id 
					FROM gachi_personas 
					WHERE team_id = :team_id 
					ORDER BY strength_level DESC`,
			Params: []model.EndpointParams{
				{
					Name:     "team_id",
					Type:     "int",
					Required: true,
				},
			},
		}
		params := map[string]any{
			"team_id": 2,
		}
		rows, err := connector.Query(ctx, endpoint, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, rows)

		// Verify results
		expected := []map[string]any{
			{"name": "Billy Herrington", "strength_level": int64(100), "special_move": "Anvil Drop", "team_id": int64(2)},
			{"name": "Hard Rod", "strength_level": int64(88), "special_move": "Steel Pipe Crush", "team_id": int64(2)},
			{"name": "Mark Wolff", "strength_level": int64(85), "special_move": "Wolf Howl Slam", "team_id": int64(2)},
			{"name": "Muscle Daddy", "strength_level": int64(79), "special_move": "Bear Hug Crush", "team_id": int64(2)},
		}

		assert.Equal(t, len(expected), len(rows))
		for i, exp := range expected {
			assert.Equal(t, exp["name"], rows[i]["name"])
			assert.Equal(t, exp["strength_level"], rows[i]["strength_level"])
			assert.Equal(t, exp["special_move"], rows[i]["special_move"])
			assert.Equal(t, exp["team_id"], rows[i]["team_id"])
		}
	})
}
