package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tucredito/backend-api/internal/domain"
)

// BankRepository implements repository.BankRepository with PostgreSQL.
type BankRepository struct {
	pool *pgxpool.Pool
}

// NewBankRepository returns a new PostgreSQL bank repository.
func NewBankRepository(pool *pgxpool.Pool) *BankRepository {
	return &BankRepository{
		pool: pool,
	}
}

// Creates a new bank and returns it.
func (r *BankRepository) Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
	id := uuid.New().String()
	var b domain.Bank
	query := `
		INSERT INTO banks (id, name, type, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, name, type
	`

	err := r.pool.QueryRow(ctx, query, id, input.Name, input.Type).Scan(&b.ID, &b.Name, &b.Type)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// Gets a bank by ID or nil if not found.
func (r *BankRepository) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	query := `SELECT id, name, type FROM banks WHERE id = $1`
	var b domain.Bank

	err := r.pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.Name, &b.Type)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &b, nil
}

// Lists banks with pagination.
func (r *BankRepository) List(ctx context.Context, limit, offset int) ([]*domain.Bank, error) {
	if limit <= 0 {
		limit = 20
	}

	var list []*domain.Bank
	query := `SELECT id, name, type FROM banks ORDER BY name LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var b domain.Bank
		if err := rows.Scan(&b.ID, &b.Name, &b.Type); err != nil {
			return nil, err
		}
		list = append(list, &b)
	}

	return list, rows.Err()
}
