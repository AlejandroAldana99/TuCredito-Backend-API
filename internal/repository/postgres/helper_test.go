package postgres_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/repository/postgres"
)

func testDBPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, connString)
	require.NoError(t, err)
	return pool
}

func deleteCredit(t *testing.T, pool *pgxpool.Pool, id string) {
	t.Helper()
	if id == "" {
		return
	}
	_, _ = pool.Exec(context.Background(), "DELETE FROM credits WHERE id = $1", id)
}

func deleteClient(t *testing.T, pool *pgxpool.Pool, id string) {
	t.Helper()
	if id == "" {
		return
	}
	_, _ = pool.Exec(context.Background(), "DELETE FROM clients WHERE id = $1", id)
}

func deleteBank(t *testing.T, pool *pgxpool.Pool, id string) {
	t.Helper()
	if id == "" {
		return
	}
	_, _ = pool.Exec(context.Background(), "DELETE FROM banks WHERE id = $1", id)
}

func uniqueClientEmail(t *testing.T) string {
	t.Helper()
	safe := strings.ReplaceAll(t.Name(), "/", "-")
	return "test-" + safe + "-" + uuid.New().String()[:8] + "@test.com"
}
