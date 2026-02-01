package service

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/repository"
)

type bankService struct {
	repository repository.BankRepository
}

func NewBankService(repository repository.BankRepository) BankService {
	return &bankService{
		repository: repository,
	}
}

// Create creates a new bank
func (s *bankService) Create(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
	return s.repository.Create(ctx, input)
}

// GetByID gets a bank by ID
func (s *bankService) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	return s.repository.GetByID(ctx, id)
}

// Lists banks with pagination
func (s *bankService) List(ctx context.Context, limit, offset int) ([]*domain.Bank, error) {
	return s.repository.List(ctx, limit, offset)
}
