package mocks

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

/*
	ClientRepository is a mock for repository.ClientRepository
	Used for testing purposes
*/

type ClientRepository struct {
	CreateFunc  func(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error)
	GetByIDFunc func(ctx context.Context, id string) (*domain.Client, error)
	ListFunc    func(ctx context.Context, limit, offset int) ([]*domain.Client, error)
}

func (m *ClientRepository) Create(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *ClientRepository) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *ClientRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}
