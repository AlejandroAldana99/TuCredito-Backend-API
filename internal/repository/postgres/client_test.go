//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
)

func TestClientRepository_Create(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewClientRepository(pool)

	input := domain.CreateClientInput{
		FullName:  "Integration Test Client",
		Email:     "integration-create@test.com",
		BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Country:   "US",
	}
	client, err := repo.Create(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.NotEmpty(t, client.ID)
	assert.Equal(t, input.FullName, client.FullName)
	assert.Equal(t, input.Email, client.Email)
	assert.Equal(t, input.Country, client.Country)
	assert.True(t, client.IsActive)
}

func TestClientRepository_GetByID(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewClientRepository(pool)

	created, err := repo.Create(ctx, domain.CreateClientInput{
		FullName:  "GetByID Client",
		Email:     "getbyid@test.com",
		BirthDate: time.Date(1985, 5, 5, 0, 0, 0, 0, time.UTC),
		Country:   "CO",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "GetByID Client", got.FullName)

	_, err = repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	// Not found returns (nil, nil)
}

func TestClientRepository_Update(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewClientRepository(pool)

	created, err := repo.Create(ctx, domain.CreateClientInput{
		FullName:  "Original Name",
		Email:     "update@test.com",
		BirthDate: time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC),
		Country:   "US",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	updated, err := repo.Update(ctx, created.ID, domain.UpdateClientInput{
		FullName:  "Updated Name",
		Email:     "updated@test.com",
		BirthDate: time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC),
		Country:   "MX",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.FullName)
	assert.Equal(t, "updated@test.com", updated.Email)
	assert.Equal(t, "MX", updated.Country)
}

func TestClientRepository_SetInactive(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewClientRepository(pool)

	created, err := repo.Create(ctx, domain.CreateClientInput{
		FullName:  "To Deactivate",
		Email:     "deactivate@test.com",
		BirthDate: time.Date(1988, 8, 8, 0, 0, 0, 0, time.UTC),
		Country:   "US",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	softDeleted, err := repo.SetInactive(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, softDeleted)
	assert.False(t, softDeleted.IsActive)
}

func TestClientRepository_List(t *testing.T) {
	pool := integrationPool(t)
	defer pool.Close()
	ctx := context.Background()
	repo := NewClientRepository(pool)

	_, err := repo.Create(ctx, domain.CreateClientInput{
		FullName:  "List Client",
		Email:     "list@test.com",
		BirthDate: time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC),
		Country:   "US",
	})
	require.NoError(t, err)

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
}
