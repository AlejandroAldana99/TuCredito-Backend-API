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
	CreateFunc      func(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error)
	GetByIDFunc     func(ctx context.Context, id string) (*domain.Client, error)
	UpdateFunc      func(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error)
	SetInactiveFunc func(ctx context.Context, id string) (*domain.Client, error)
	SetActiveFunc   func(ctx context.Context, id string) (*domain.Client, error)
	ListFunc        func(ctx context.Context, limit, offset int) ([]*domain.Client, error)
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

func (m *ClientRepository) Update(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *ClientRepository) SetInactive(ctx context.Context, id string) (*domain.Client, error) {
	if m.SetInactiveFunc != nil {
		return m.SetInactiveFunc(ctx, id)
	}
	return nil, nil
}

func (m *ClientRepository) SetActive(ctx context.Context, id string) (*domain.Client, error) {
	if m.SetActiveFunc != nil {
		return m.SetActiveFunc(ctx, id)
	}
	return nil, nil
}

func (m *ClientRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}
