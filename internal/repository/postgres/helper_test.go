//go:build integration
// +build integration

package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func integrationPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}
	ctx := context.Background()
	pool, err := NewPool(ctx, connString)
	require.NoError(t, err)
	return pool
}
