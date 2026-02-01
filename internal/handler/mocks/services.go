package mocks

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/service"
)

type MockClientService struct {
	CreateFunc  func(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error)
	GetByIDFunc func(ctx context.Context, id string) (*domain.Client, error)
	UpdateFunc  func(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error)
	DeleteFunc  func(ctx context.Context, id string) (*domain.Client, error)
	ListFunc    func(ctx context.Context, limit, offset int) ([]*domain.Client, error)
}

func (m *MockClientService) Create(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockClientService) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClientService) Update(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *MockClientService) Delete(ctx context.Context, id string) (*domain.Client, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClientService) List(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}

var _ service.ClientService = (*MockClientService)(nil)

type MockBankService struct {
	CreateFunc  func(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error)
	GetByIDFunc func(ctx context.Context, id string) (*domain.Bank, error)
	UpdateFunc  func(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error)
	DeleteFunc  func(ctx context.Context, id string) (*domain.Bank, error)
	ListFunc    func(ctx context.Context, limit, offset int) ([]*domain.Bank, error)
}

func (m *MockBankService) Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockBankService) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBankService) Update(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *MockBankService) Delete(ctx context.Context, id string) (*domain.Bank, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBankService) List(ctx context.Context, limit, offset int) ([]*domain.Bank, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}

var _ service.BankService = (*MockBankService)(nil)

type MockCreditService struct {
	CreateFunc         func(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error)
	GetByIDFunc        func(ctx context.Context, id string) (*domain.Credit, error)
	UpdateFunc         func(ctx context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error)
	UpdateStatusFunc   func(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error)
	DeleteFunc         func(ctx context.Context, id string) (*domain.Credit, error)
	ListFunc           func(ctx context.Context, limit, offset int) ([]*domain.Credit, error)
	ListByClientIDFunc func(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error)
}

func (m *MockCreditService) Create(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockCreditService) GetByID(ctx context.Context, id string) (*domain.Credit, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCreditService) Update(ctx context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *MockCreditService) UpdateStatus(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error) {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil, nil
}

func (m *MockCreditService) Delete(ctx context.Context, id string) (*domain.Credit, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCreditService) List(ctx context.Context, limit, offset int) ([]*domain.Credit, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockCreditService) ListByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error) {
	if m.ListByClientIDFunc != nil {
		return m.ListByClientIDFunc(ctx, clientID, limit, offset)
	}
	return nil, nil
}

func (m *MockCreditService) Shutdown() {}

var _ service.CreditService = (*MockCreditService)(nil)
