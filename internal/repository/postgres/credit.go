package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tucredito/backend-api/internal/domain"
)

type CreditRepository struct {
	pool *pgxpool.Pool
}

// Creates a new PostgreSQL credit repository
func NewCreditRepository(pool *pgxpool.Pool) *CreditRepository {
	return &CreditRepository{pool: pool}
}

// Inserts a new credit and returns it
func (r *CreditRepository) Create(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO credits (id, client_id, bank_id, min_payment, max_payment, term_months, credit_type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING', NOW(), NOW())
		RETURNING id, client_id, bank_id, min_payment, max_payment, term_months, credit_type, status, created_at
	`

	var c domain.Credit
	err := r.pool.QueryRow(ctx, query, id, input.ClientID, input.BankID, input.MinPayment, input.MaxPayment, input.TermMonths, input.CreditType).Scan(
		&c.ID, &c.ClientID, &c.BankID, &c.MinPayment, &c.MaxPayment, &c.TermMonths, &c.CreditType, &c.Status, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Returns a credit by ID or nil if not found
func (r *CreditRepository) GetByID(ctx context.Context, id string) (*domain.Credit, error) {
	var c domain.Credit
	query := `
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, credit_type, status, created_at
		FROM credits WHERE id = $1
	`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.ClientID, &c.BankID, &c.MinPayment, &c.MaxPayment,
		&c.TermMonths, &c.CreditType, &c.Status, &c.CreatedAt,
	)

	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

// Updates credit status and returns the credit
func (r *CreditRepository) UpdateStatus(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error) {
	query := `
		UPDATE credits SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, client_id, bank_id, min_payment, max_payment, term_months, credit_type, status, created_at
	`

	var c domain.Credit
	err := r.pool.QueryRow(ctx, query, status, id).Scan(
		&c.ID, &c.ClientID, &c.BankID, &c.MinPayment, &c.MaxPayment,
		&c.TermMonths, &c.CreditType, &c.Status, &c.CreatedAt,
	)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

// List returns credits with pagination
func (r *CreditRepository) List(ctx context.Context, limit, offset int) ([]*domain.Credit, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, credit_type, status, created_at
		FROM credits ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanCredits(rows)
}

// Returns credits for a client with pagination
func (r *CreditRepository) ListByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, client_id, bank_id, min_payment, max_payment, term_months, credit_type, status, created_at
		FROM credits WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanCredits(rows)
}
