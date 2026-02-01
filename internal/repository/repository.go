package repository

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

// ClientRepository defines the methods for client repository persistence
type ClientRepository interface {
	Create(ctx context.Context, client domain.CreateClientInput) (*domain.Client, error)
	GetByID(ctx context.Context, id string) (*domain.Client, error)
	Update(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error)
	SetInactive(ctx context.Context, id string) (*domain.Client, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Client, error)
}

// BankRepository defines the methods for bank repository persistence
type BankRepository interface {
	Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error)
	GetByID(ctx context.Context, id string) (*domain.Bank, error)
	Update(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error)
	SetInactive(ctx context.Context, id string) (*domain.Bank, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Bank, error)
}

// CreditRepository defines the methods for credit repository persistence
type CreditRepository interface {
	Create(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error)
	GetByID(ctx context.Context, id string) (*domain.Credit, error)
	Update(ctx context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error)
	UpdateStatus(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error)
	SetInactive(ctx context.Context, id string) (*domain.Credit, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Credit, error)
	ListByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error)
}
