package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
	handlermocks "github.com/tucredito/backend-api/internal/handler/mocks"
	"go.uber.org/zap"
)

const apiVersion = "/v1"

func TestClientHandler_Create(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	created := &domain.Client{
		ID: "client-1", FullName: "Jane Doe", Email: "jane@example.com",
		Country: "US", IsActive: true, CreatedAt: time.Now(),
	}
	mockSvc.CreateFunc = func(_ context.Context, input domain.CreateClientInput) (*domain.Client, error) {
		out := *created
		out.FullName = input.FullName
		out.Email = input.Email
		out.Country = input.Country
		return &out, nil
	}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/clients", h.Create)

	body := []byte(`{"full_name":"Jane Doe","email":"jane@example.com","country":"US"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/clients", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got domain.Client
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "Jane Doe", got.FullName)
	assert.Equal(t, "jane@example.com", got.Email)
	assert.Equal(t, "US", got.Country)
}

func TestClientHandler_Create_ValidationError(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/clients", h.Create)

	body := []byte(`{"full_name":"","email":"jane@example.com","country":"US"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/clients", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errBody struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&errBody))
	assert.Equal(t, "VALIDATION", errBody.Code)
}

func TestClientHandler_Create_InvalidJSON(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/clients", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/v1/clients", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestClientHandler_GetByID(t *testing.T) {
	log, _ := zap.NewDevelopment()
	client := &domain.Client{ID: "c1", FullName: "Jane", Email: "j@x.com", Country: "US", IsActive: true}
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.GetByIDFunc = func(_ context.Context, id string) (*domain.Client, error) {
		if id == "c1" {
			return client, nil
		}
		return nil, nil
	}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/clients/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/clients/c1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Client
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "c1", got.ID)
	assert.Equal(t, "Jane", got.FullName)
}

func TestClientHandler_GetByID_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.GetByIDFunc = func(_ context.Context, _ string) (*domain.Client, error) { return nil, nil }
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/clients/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/clients/none", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestClientHandler_List(t *testing.T) {
	log, _ := zap.NewDevelopment()
	list := []*domain.Client{
		{ID: "c1", FullName: "A", Email: "a@b.com", Country: "US", IsActive: true},
	}
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.ListFunc = func(_ context.Context, limit, offset int) ([]*domain.Client, error) { return list, nil }
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/clients", h.List)

	req := httptest.NewRequest(http.MethodGet, "/v1/clients", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []*domain.Client
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	require.Len(t, got, 1)
	assert.Equal(t, "c1", got[0].ID)
}

func TestClientHandler_Update(t *testing.T) {
	log, _ := zap.NewDevelopment()
	updated := &domain.Client{ID: "c1", FullName: "Jane Updated", Email: "j2@x.com", Country: "US", IsActive: true}
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.UpdateFunc = func(_ context.Context, id string, input domain.UpdateClientInput) (*domain.Client, error) {
		out := *updated
		out.FullName = input.FullName
		out.Email = input.Email
		out.Country = input.Country
		return &out, nil
	}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT "+apiVersion+"/clients/{id}", h.Update)

	body := []byte(`{"full_name":"Jane Updated","email":"j2@x.com","country":"US"}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/clients/c1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Client
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "Jane Updated", got.FullName)
}

func TestClientHandler_Update_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.UpdateFunc = func(_ context.Context, _ string, _ domain.UpdateClientInput) (*domain.Client, error) { return nil, nil }
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT "+apiVersion+"/clients/{id}", h.Update)

	body := []byte(`{"full_name":"X","email":"x@x.com","country":"US"}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/clients/none", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestClientHandler_Delete(t *testing.T) {
	log, _ := zap.NewDevelopment()
	softDeleted := &domain.Client{ID: "c1", FullName: "Jane", IsActive: false}
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.DeleteFunc = func(_ context.Context, id string) (*domain.Client, error) {
		if id == "c1" {
			return softDeleted, nil
		}
		return nil, nil
	}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE "+apiVersion+"/clients/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/v1/clients/c1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestClientHandler_Delete_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.DeleteFunc = func(_ context.Context, _ string) (*domain.Client, error) { return nil, nil }
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE "+apiVersion+"/clients/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/v1/clients/none", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestClientHandler_Reenable(t *testing.T) {
	log, _ := zap.NewDevelopment()
	reenabled := &domain.Client{ID: "c1", FullName: "Jane", IsActive: true}
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.ReenableFunc = func(_ context.Context, id string) (*domain.Client, error) {
		if id == "c1" {
			return reenabled, nil
		}
		return nil, nil
	}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/clients/{id}/reenable", h.Reenable)

	req := httptest.NewRequest(http.MethodPost, "/v1/clients/c1/reenable", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Client
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "c1", got.ID)
	assert.True(t, got.IsActive)
}

func TestClientHandler_Reenable_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	mockSvc.ReenableFunc = func(_ context.Context, _ string) (*domain.Client, error) { return nil, nil }
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/clients/{id}/reenable", h.Reenable)

	req := httptest.NewRequest(http.MethodPost, "/v1/clients/none/reenable", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestClientHandler_MethodNotAllowed(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockClientService{}
	h := NewClientHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/clients", h.Create)

	req := httptest.NewRequest(http.MethodGet, "/v1/clients", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
