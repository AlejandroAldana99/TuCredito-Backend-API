package repository

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

// ClientRepository defines the persistence for clients.
type ClientRepository interface {
	Create(ctx context.Context, client domain.CreateClientInput) (*domain.Client, error)
	GetByID(ctx context.Context, id string) (*domain.Client, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Client, error)
}

// BankRepository defines persistence operations for banks.
type BankRepository interface {
	Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error)
	GetByID(ctx context.Context, id string) (*domain.Bank, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Bank, error)
}

// CreditRepository defines persistence operations for credits.
type CreditRepository interface {
	Create(ctx context.Context, input domain.CreateCreditInput) (*domain.Credit, error)
	GetByID(ctx context.Context, id string) (*domain.Credit, error)
	UpdateStatus(ctx context.Context, id string, status domain.CreditStatus) (*domain.Credit, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Credit, error)
	ListByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Credit, error)
}
