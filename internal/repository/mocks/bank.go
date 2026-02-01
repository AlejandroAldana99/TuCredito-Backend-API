package mocks

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

/*
	BankRepository is a mock for repository.BankRepository
	Used for testing purposes
*/

type BankRepository struct {
	CreateFunc      func(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error)
	GetByIDFunc     func(ctx context.Context, id string) (*domain.Bank, error)
	UpdateFunc      func(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error)
	SetInactiveFunc func(ctx context.Context, id string) (*domain.Bank, error)
	ListFunc        func(ctx context.Context, limit, offset int) ([]*domain.Bank, error)
}

func (m *BankRepository) Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *BankRepository) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *BankRepository) Update(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *BankRepository) SetInactive(ctx context.Context, id string) (*domain.Bank, error) {
	if m.SetInactiveFunc != nil {
		return m.SetInactiveFunc(ctx, id)
	}
	return nil, nil
}

func (m *BankRepository) List(ctx context.Context, limit, offset int) ([]*domain.Bank, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}
