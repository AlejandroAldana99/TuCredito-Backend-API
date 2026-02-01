//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
)

func TestBankRepository_Create(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewBankRepository(pool)

	input := domain.CreateBankInput{Name: "Integration Test Bank", Type: domain.BankTypePrivate}
	bank, err := repo.Create(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, bank)
	assert.NotEmpty(t, bank.ID)
	assert.Equal(t, input.Name, bank.Name)
	assert.Equal(t, input.Type, bank.Type)
	assert.True(t, bank.IsActive)
}

func TestBankRepository_GetByID(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewBankRepository(pool)

	created, err := repo.Create(ctx, domain.CreateBankInput{Name: "GetByID Bank", Type: domain.BankTypeGovernment})
	require.NoError(t, err)
	require.NotNil(t, created)

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "GetByID Bank", got.Name)
	assert.Equal(t, domain.BankTypeGovernment, got.Type)
}

func TestBankRepository_Update(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewBankRepository(pool)

	created, err := repo.Create(ctx, domain.CreateBankInput{Name: "Original Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, created)

	updated, err := repo.Update(ctx, created.ID, domain.UpdateBankInput{Name: "Updated Bank", Type: domain.BankTypeGovernment})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated Bank", updated.Name)
	assert.Equal(t, domain.BankTypeGovernment, updated.Type)
}

func TestBankRepository_SetInactive(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewBankRepository(pool)

	created, err := repo.Create(ctx, domain.CreateBankInput{Name: "To Deactivate Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, created)

	softDeleted, err := repo.SetInactive(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, softDeleted)
	assert.False(t, softDeleted.IsActive)
}

func TestBankRepository_List(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewBankRepository(pool)

	_, err := repo.Create(ctx, domain.CreateBankInput{Name: "List Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
}
