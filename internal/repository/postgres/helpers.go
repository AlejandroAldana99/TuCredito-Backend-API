package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tucredito/backend-api/internal/domain"
)

// Creates a PostgreSQL connection pool
func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, config)
}

// Checks if the error is "no rows"
func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// scanCredits scans credit rows into a slice
func scanCredits(rows pgx.Rows) ([]*domain.Credit, error) {
	var list []*domain.Credit
	for rows.Next() {
		var c domain.Credit
		if err := rows.Scan(
			&c.ID, &c.ClientID, &c.BankID, &c.MinPayment, &c.MaxPayment,
			&c.TermMonths, &c.CreditType, &c.Status, &c.CreatedAt, &c.IsActive,
		); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}
