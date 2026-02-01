package service

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

// ClientService defines the methods for client service logic
type ClientService interface {
	Create(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error)
	GetByID(ctx context.Context, id string) (*domain.Client, error)
	Update(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error)
	Delete(ctx context.Context, id string) (*domain.Client, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Client, error)
}

// BankService defines the methods for bank service logic
type BankService interface {
	Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error)
	GetByID(ctx context.Context, id string) (*domain.Bank, error)
	Update(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error)
	Delete(ctx context.Context, id string) (*domain.Bank, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Bank, error)
}

// CreditService defines the methods for credit service logic
type CreditService interface {
	Create(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error)
	GetByID(ctx context.Context, id string) (*domain.Credit, error)
	Update(ctx context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error)
	UpdateStatus(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error)
	Delete(ctx context.Context, id string) (*domain.Credit, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Credit, error)
	ListByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error)
	Shutdown()
}
