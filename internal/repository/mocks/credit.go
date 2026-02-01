package mocks

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

/*
	CreditRepository is a mock for repository.CreditRepository
	Used for testing purposes
*/

type CreditRepository struct {
	CreateFunc         func(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error)
	GetByIDFunc        func(ctx context.Context, id string) (*domain.Credit, error)
	UpdateFunc         func(ctx context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error)
	UpdateStatusFunc   func(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error)
	SetInactiveFunc    func(ctx context.Context, id string) (*domain.Credit, error)
	SetActiveFunc      func(ctx context.Context, id string) (*domain.Credit, error)
	ListFunc           func(ctx context.Context, limit, offset int) ([]*domain.Credit, error)
	ListByClientIDFunc func(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error)
}

func (m *CreditRepository) Create(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *CreditRepository) GetByID(ctx context.Context, id string) (*domain.Credit, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *CreditRepository) Update(ctx context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *CreditRepository) UpdateStatus(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error) {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil, nil
}

func (m *CreditRepository) SetInactive(ctx context.Context, id string) (*domain.Credit, error) {
	if m.SetInactiveFunc != nil {
		return m.SetInactiveFunc(ctx, id)
	}
	return nil, nil
}

func (m *CreditRepository) SetActive(ctx context.Context, id string) (*domain.Credit, error) {
	if m.SetActiveFunc != nil {
		return m.SetActiveFunc(ctx, id)
	}
	return nil, nil
}

func (m *CreditRepository) List(ctx context.Context, limit, offset int) ([]*domain.Credit, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *CreditRepository) ListByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error) {
	if m.ListByClientIDFunc != nil {
		return m.ListByClientIDFunc(ctx, clientID, limit, offset)
	}
	return nil, nil
}
