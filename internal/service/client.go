package service

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/repository"
)

type clientService struct {
	repository repository.ClientRepository
}

func NewClientService(repository repository.ClientRepository) ClientService {
	return &clientService{
		repository: repository,
	}
}

// Create creates a new client
func (s *clientService) Create(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error) {
	return s.repository.Create(ctx, input)
}

// GetByID gets a client by ID
func (s *clientService) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	return s.repository.GetByID(ctx, id)
}

// Lists clients with pagination
func (s *clientService) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	return s.repository.List(ctx, limit, offset)
}
