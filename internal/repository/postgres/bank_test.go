package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/repository/postgres"
)

func TestBankRepository_Create(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := postgres.NewBankRepository(pool)

	input := domain.CreateBankInput{Name: "Integration Test Bank", Type: domain.BankTypePrivate}
	bank, err := repo.Create(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)
	assert.NotEmpty(t, bank.ID)
	assert.Equal(t, input.Name, bank.Name)
	assert.Equal(t, input.Type, bank.Type)
	assert.True(t, bank.IsActive)
}

func TestBankRepository_GetByID(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := postgres.NewBankRepository(pool)

	created, err := repo.Create(ctx, domain.CreateBankInput{Name: "GetByID Bank", Type: domain.BankTypeGovernment})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteBank(t, pool, created.ID)

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "GetByID Bank", got.Name)
	assert.Equal(t, domain.BankTypeGovernment, got.Type)
}

func TestBankRepository_Update(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := postgres.NewBankRepository(pool)

	created, err := repo.Create(ctx, domain.CreateBankInput{Name: "Original Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteBank(t, pool, created.ID)

	updated, err := repo.Update(ctx, created.ID, domain.UpdateBankInput{Name: "Updated Bank", Type: domain.BankTypeGovernment})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated Bank", updated.Name)
	assert.Equal(t, domain.BankTypeGovernment, updated.Type)
}

func TestBankRepository_SetInactive(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := postgres.NewBankRepository(pool)

	created, err := repo.Create(ctx, domain.CreateBankInput{Name: "To Deactivate Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteBank(t, pool, created.ID)

	softDeleted, err := repo.SetInactive(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, softDeleted)
	assert.False(t, softDeleted.IsActive)
}

func TestBankRepository_SetActive(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := postgres.NewBankRepository(pool)

	created, err := repo.Create(ctx, domain.CreateBankInput{Name: "To Reenable Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, created)
	defer deleteBank(t, pool, created.ID)

	_, err = repo.SetInactive(ctx, created.ID)
	require.NoError(t, err)

	reenabled, err := repo.SetActive(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, reenabled)
	assert.True(t, reenabled.IsActive)
	assert.Equal(t, created.ID, reenabled.ID)
}

func TestBankRepository_SetActive_NotFound(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := postgres.NewBankRepository(pool)

	got, err := repo.SetActive(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestBankRepository_List(t *testing.T) {
	pool := testDBPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := postgres.NewBankRepository(pool)

	bank, err := repo.Create(ctx, domain.CreateBankInput{Name: "List Bank", Type: domain.BankTypePrivate})
	require.NoError(t, err)
	require.NotNil(t, bank)
	defer deleteBank(t, pool, bank.ID)

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
}
