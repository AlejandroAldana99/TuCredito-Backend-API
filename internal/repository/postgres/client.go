package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tucredito/backend-api/internal/domain"
)

// ClientRepository implements ClientRepository with PostgreSQL.
type ClientRepository struct {
	pool *pgxpool.Pool
}

// NewClientRepository returns a new ClientRepository with a PostgreSQL pool.
func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{
		pool: pool,
	}
}

// Inserts a new client and returns it.
func (r *ClientRepository) Create(ctx context.Context, client domain.CreateClientInput) (*domain.Client, error) {
	// Generate a new UUID for the client
	id := uuid.New().String()
	var c domain.Client
	query := `
		INSERT INTO clients (id, full_name, email, birth_date, country, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, full_name, email, birth_date, country, created_at
	`
	// Execute the query and scan the result
	err := r.pool.QueryRow(ctx, query,
		id, client.FullName, client.Email, client.BirthDate, client.Country,
	).Scan(&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Gets a client by ID or nil if not found.
func (r *ClientRepository) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	var c domain.Client
	query := `
		SELECT id, full_name, email, birth_date, country, created_at
		FROM clients WHERE id = $1
	`

	// Execute the query and scan the result
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt,
	)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

// Lists clients with pagination.
func (r *ClientRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	// If the limit is less than or equal to 0, set it to 20
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, full_name, email, birth_date, country, created_at
		FROM clients ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`

	// Execute the query and get the rows
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(&c.ID, &c.FullName, &c.Email, &c.BirthDate, &c.Country, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}

	return list, rows.Err()
}
