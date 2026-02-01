package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
	repomocks "github.com/tucredito/backend-api/internal/repository/mocks"
)

func TestClientService_Create(t *testing.T) {
	created := &domain.Client{
		ID: "c1", FullName: "Jane Doe", Email: "jane@example.com",
		Country: "US", IsActive: true, CreatedAt: time.Now(),
	}
	repo := &repomocks.ClientRepository{}
	repo.CreateFunc = func(ctx context.Context, input domain.CreateClientInput) (*domain.Client, error) {
		out := *created
		out.FullName = input.FullName
		out.Email = input.Email
		out.Country = input.Country
		return &out, nil
	}
	svc := NewClientService(repo)

	got, err := svc.Create(context.Background(), domain.CreateClientInput{
		FullName: "Jane Doe", Email: "jane@example.com", Country: "US",
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "c1", got.ID)
	assert.Equal(t, "Jane Doe", got.FullName)
	assert.Equal(t, "jane@example.com", got.Email)
}

func TestClientService_GetByID_Found(t *testing.T) {
	client := &domain.Client{ID: "c1", FullName: "Jane", Email: "j@x.com", Country: "US", IsActive: true}
	repo := &repomocks.ClientRepository{}
	repo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Client, error) {
		if id == "c1" {
			return client, nil
		}
		return nil, nil
	}
	svc := NewClientService(repo)

	got, err := svc.GetByID(context.Background(), "c1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "c1", got.ID)
}

func TestClientService_GetByID_NotFound(t *testing.T) {
	repo := &repomocks.ClientRepository{}
	repo.GetByIDFunc = func(ctx context.Context, id string) (*domain.Client, error) {
		return nil, nil
	}
	svc := NewClientService(repo)

	got, err := svc.GetByID(context.Background(), "none")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestClientService_Update(t *testing.T) {
	updated := &domain.Client{ID: "c1", FullName: "Jane Updated", Email: "j2@x.com", Country: "US", IsActive: true}
	repo := &repomocks.ClientRepository{}
	repo.UpdateFunc = func(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error) {
		out := *updated
		out.FullName = input.FullName
		out.Email = input.Email
		out.Country = input.Country
		return &out, nil
	}
	svc := NewClientService(repo)

	got, err := svc.Update(context.Background(), "c1", domain.UpdateClientInput{
		FullName: "Jane Updated", Email: "j2@x.com", Country: "US",
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Jane Updated", got.FullName)
}

func TestClientService_Update_NotFound(t *testing.T) {
	repo := &repomocks.ClientRepository{}
	repo.UpdateFunc = func(ctx context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error) {
		return nil, nil
	}
	svc := NewClientService(repo)

	got, err := svc.Update(context.Background(), "none", domain.UpdateClientInput{
		FullName: "X", Email: "x@x.com", Country: "US",
	})
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestClientService_Delete(t *testing.T) {
	softDeleted := &domain.Client{ID: "c1", FullName: "Jane", IsActive: false}
	repo := &repomocks.ClientRepository{}
	repo.SetInactiveFunc = func(ctx context.Context, id string) (*domain.Client, error) {
		if id == "c1" {
			return softDeleted, nil
		}
		return nil, nil
	}
	svc := NewClientService(repo)

	got, err := svc.Delete(context.Background(), "c1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.False(t, got.IsActive)
}

func TestClientService_List(t *testing.T) {
	list := []*domain.Client{
		{ID: "c1", FullName: "A", Email: "a@b.com", Country: "US", IsActive: true},
	}
	repo := &repomocks.ClientRepository{}
	repo.ListFunc = func(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
		return list, nil
	}
	svc := NewClientService(repo)

	got, err := svc.List(context.Background(), 10, 0)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "c1", got[0].ID)
}
