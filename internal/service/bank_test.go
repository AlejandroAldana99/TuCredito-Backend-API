package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
	repomocks "github.com/tucredito/backend-api/internal/repository/mocks"
)

func TestBankService_Create(t *testing.T) {
	created := &domain.Bank{ID: "b1", Name: "Test Bank", Type: domain.BankTypePrivate, IsActive: true}
	repo := &repomocks.BankRepository{}
	repo.CreateFunc = func(ctx context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
		out := *created
		out.Name = input.Name
		out.Type = input.Type
		return &out, nil
	}
	svc := NewBankService(repo)

	got, err := svc.Create(context.Background(), domain.CreateBankInput{
		Name: "Test Bank", Type: domain.BankTypePrivate,
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "b1", got.ID)
	assert.Equal(t, "Test Bank", got.Name)
	assert.Equal(t, domain.BankTypePrivate, got.Type)
}

func TestBankService_GetByID_Found(t *testing.T) {
	bank := &domain.Bank{ID: "b1", Name: "Bank One", Type: domain.BankTypeGovernment, IsActive: true}
	repo := &repomocks.BankRepository{}
	repo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Bank, error) {
		if id == "b1" {
			return bank, nil
		}
		return nil, nil
	}
	svc := NewBankService(repo)

	got, err := svc.GetByID(context.Background(), "b1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "b1", got.ID)
	assert.Equal(t, domain.BankTypeGovernment, got.Type)
}

func TestBankService_GetByID_NotFound(t *testing.T) {
	repo := &repomocks.BankRepository{}
	repo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Bank, error) {
		return nil, nil
	}
	svc := NewBankService(repo)

	got, err := svc.GetByID(context.Background(), "none")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestBankService_Update(t *testing.T) {
	updated := &domain.Bank{ID: "b1", Name: "Bank Updated", Type: domain.BankTypePrivate, IsActive: true}
	repo := &repomocks.BankRepository{}
	repo.UpdateFunc = func(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error) {
		out := *updated
		out.Name = input.Name
		out.Type = input.Type
		return &out, nil
	}
	svc := NewBankService(repo)

	got, err := svc.Update(context.Background(), "b1", domain.UpdateBankInput{
		Name: "Bank Updated", Type: domain.BankTypePrivate,
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Bank Updated", got.Name)
}

func TestBankService_Update_NotFound(t *testing.T) {
	repo := &repomocks.BankRepository{}
	repo.UpdateFunc = func(ctx context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error) {
		return nil, nil
	}
	svc := NewBankService(repo)

	got, err := svc.Update(context.Background(), "none", domain.UpdateBankInput{
		Name: "X", Type: domain.BankTypePrivate,
	})
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestBankService_Delete(t *testing.T) {
	softDeleted := &domain.Bank{ID: "b1", Name: "Bank", IsActive: false}
	repo := &repomocks.BankRepository{}
	repo.SetInactiveFunc = func(ctx context.Context, id string) (*domain.Bank, error) {
		if id == "b1" {
			return softDeleted, nil
		}
		return nil, nil
	}
	svc := NewBankService(repo)

	got, err := svc.Delete(context.Background(), "b1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.False(t, got.IsActive)
}

func TestBankService_List(t *testing.T) {
	list := []*domain.Bank{
		{ID: "b1", Name: "Bank A", Type: domain.BankTypePrivate, IsActive: true},
	}
	repo := &repomocks.BankRepository{}
	repo.ListFunc = func(ctx context.Context, limit, offset int) ([]*domain.Bank, error) {
		return list, nil
	}
	svc := NewBankService(repo)

	got, err := svc.List(context.Background(), 10, 0)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "b1", got[0].ID)
}
