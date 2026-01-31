package service

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/repository"
)

// ClientService handles client business logic.
type ClientService struct {
	repository repository.ClientRepository
}

// NewClientService returns a new ClientService.
func NewClientService(repository repository.ClientRepository) *ClientService {
	return &ClientService{
		repository: repository,
	}
}

// Create a new client.
func (s *ClientService) Create(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error) {
	// Create the client using the repository
	return s.repository.Create(ctx, input)
}

// Get a client by ID.
func (s *ClientService) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	// Get the client by ID using the repository
	return s.repository.GetByID(ctx, id)
}

// List clients with pagination.
func (s *ClientService) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	// List the clients using the repository
	return s.repository.List(ctx, limit, offset)
}
