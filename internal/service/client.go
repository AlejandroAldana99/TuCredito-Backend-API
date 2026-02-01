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

// Creates a new client
func (s *clientService) Create(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error) {
	return s.repository.Create(ctx, input)
}

// Gets a client by ID
func (s *clientService) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	return s.repository.GetByID(ctx, id)
}

// Updates a client
func (s *clientService) Update(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error) {
	return s.repository.Update(ctx, id, input)
}

// Soft-deletes a client
func (s *clientService) Delete(ctx context.Context, id string) (*domain.Client, error) {
	return s.repository.SetInactive(ctx, id)
}

// Lists clients with pagination
func (s *clientService) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	return s.repository.List(ctx, limit, offset)
}
