package domain

import "time"

// Client structure
type Client struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// Structure for creating a client
type CreateClientInput struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	Country   string    `json:"country"`
}

// Structure for updating a client
type UpdateClientInput struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	Country   string    `json:"country"`
}
