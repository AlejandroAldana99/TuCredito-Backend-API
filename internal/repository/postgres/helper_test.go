package postgres_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tucredito/backend-api/internal/repository/postgres"
)

func testDBPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	// Skip integration tests when DB is not available (e.g. CI test job, deploy without DB).
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
		return nil
	}
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, connString)
	if err != nil {
		t.Skipf("skipping integration test: database not available: %v", err)
		return nil
	}
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
