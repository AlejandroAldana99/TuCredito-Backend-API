package service

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/repository"
)

// BankService handles bank business logic.
type BankService struct {
	repository repository.BankRepository
}

// NewBankService returns a new BankService.
func NewBankService(repository repository.BankRepository) *BankService {
	return &BankService{
		repository: repository,
	}
}

// Creates a new bank.
func (s *BankService) Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
	return s.repository.Create(ctx, input)
}

// Gets a bank by ID.
func (s *BankService) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	return s.repository.GetByID(ctx, id)
}

// Lists banks with pagination.
func (s *BankService) List(ctx context.Context, limit, offset int) ([]*domain.Bank, error) {
	return s.repository.List(ctx, limit, offset)
}
