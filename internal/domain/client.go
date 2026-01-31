package domain

import "time"

// Client represents a client in the system.
type Client struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateClientInput holds the data required to create a client.
type CreateClientInput struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	Country   string    `json:"country"`
}
