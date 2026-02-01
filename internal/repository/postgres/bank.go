package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tucredito/backend-api/internal/domain"
)

type BankRepository struct {
	pool *pgxpool.Pool
}

func NewBankRepository(pool *pgxpool.Pool) *BankRepository {
	return &BankRepository{
		pool: pool,
	}
}

// Creates a new bank
func (r *BankRepository) Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
	id := uuid.New().String()
	var b domain.Bank
	query := `
		INSERT INTO banks (id, name, type, created_at, is_active)
		VALUES ($1, $2, $3, NOW(), TRUE)
		RETURNING id, name, type, is_active
	`
	err := r.pool.QueryRow(ctx, query, id, input.Name, input.Type).Scan(&b.ID, &b.Name, &b.Type, &b.IsActive)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Gets a bank by ID
func (r *BankRepository) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	query := `SELECT id, name, type, is_active FROM banks WHERE id = $1 AND is_active = TRUE`
	var b domain.Bank
	err := r.pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.Name, &b.Type, &b.IsActive)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// Updates a bank
func (r *BankRepository) Update(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error) {
	query := `UPDATE banks SET name = $1, type = $2 WHERE id = $3 RETURNING id, name, type, is_active`
	var b domain.Bank
	err := r.pool.QueryRow(ctx, query, input.Name, input.Type, id).Scan(&b.ID, &b.Name, &b.Type, &b.IsActive)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// Soft-deletes a bank
func (r *BankRepository) SetInactive(ctx context.Context, id string) (*domain.Bank, error) {
	query := `UPDATE banks SET is_active = FALSE WHERE id = $1 RETURNING id, name, type, is_active`
	var b domain.Bank
	err := r.pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.Name, &b.Type, &b.IsActive)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// Re-enables a bank
func (r *BankRepository) SetActive(ctx context.Context, id string) (*domain.Bank, error) {
	query := `UPDATE banks SET is_active = TRUE WHERE id = $1 RETURNING id, name, type, is_active`
	var b domain.Bank
	err := r.pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.Name, &b.Type, &b.IsActive)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// Lists banks with pagination
func (r *BankRepository) List(ctx context.Context, limit, offset int) ([]*domain.Bank, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `SELECT id, name, type, is_active FROM banks WHERE is_active = TRUE ORDER BY name LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Bank
	for rows.Next() {
		var b domain.Bank
		if err := rows.Scan(&b.ID, &b.Name, &b.Type, &b.IsActive); err != nil {
			return nil, err
		}
		list = append(list, &b)
	}
	return list, rows.Err()
}
