package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tucredito/backend-api/internal/domain"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{
		pool: pool,
	}
}

// Creates a new client
func (r *ClientRepository) Create(ctx context.Context, client domain.CreateClientInput) (*domain.Client, error) {
	// Generate a new UUID for the client
	id := uuid.New().String()
	var c domain.Client
	query := `
		INSERT INTO clients (id, full_name, email, birth_date, country, created_at, is_active)
		VALUES ($1, $2, $3, $4, $5, NOW(), TRUE)
		RETURNING id, full_name, email, birth_date, country, created_at, is_active
	`
	err := r.pool.QueryRow(ctx, query,
		id, client.FullName, client.Email, client.BirthDate, client.Country,
	).Scan(&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt, &c.IsActive)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Gets a client by ID
func (r *ClientRepository) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	var c domain.Client
	query := `SELECT id, full_name, email, birth_date, country, created_at, is_active FROM clients WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt, &c.IsActive,
	)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// Updates a client
func (r *ClientRepository) Update(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error) {
	query := `
		UPDATE clients SET full_name = $1, email = $2, birth_date = $3, country = $4
		WHERE id = $5
		RETURNING id, full_name, email, birth_date, country, created_at, is_active
	`
	var c domain.Client
	err := r.pool.QueryRow(ctx, query, input.FullName, input.Email, input.BirthDate, input.Country, id).Scan(
		&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt, &c.IsActive,
	)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// Soft-deletes a client
func (r *ClientRepository) SetInactive(ctx context.Context, id string) (*domain.Client, error) {
	query := `
		UPDATE clients SET is_active = FALSE WHERE id = $1
		RETURNING id, full_name, email, birth_date, country, created_at, is_active
	`
	var c domain.Client
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt, &c.IsActive,
	)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// Lists clients with pagination
func (r *ClientRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, full_name, email, birth_date, country, created_at, is_active
		FROM clients WHERE is_active = TRUE ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt, &c.IsActive); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}
