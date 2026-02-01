package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
	handlermocks "github.com/tucredito/backend-api/internal/handler/mocks"
	"go.uber.org/zap"
)

func TestBankHandler_Create(t *testing.T) {
	log, _ := zap.NewDevelopment()
	created := &domain.Bank{ID: "bank-1", Name: "Test Bank", Type: domain.BankTypePrivate, IsActive: true}
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.CreateFunc = func(_ context.Context, input domain.CreateBankInput) (*domain.Bank, error) {
		out := *created
		out.Name = input.Name
		out.Type = input.Type
		return &out, nil
	}
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/banks", h.Create)

	body := []byte(`{"name":"Test Bank","type":"PRIVATE"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/banks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got domain.Bank
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "Test Bank", got.Name)
	assert.Equal(t, domain.BankTypePrivate, got.Type)
}

func TestBankHandler_Create_ValidationError(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockBankService{}
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/banks", h.Create)

	body := []byte(`{"name":"","type":"PRIVATE"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/banks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBankHandler_GetByID(t *testing.T) {
	log, _ := zap.NewDevelopment()
	bank := &domain.Bank{ID: "b1", Name: "Bank One", Type: domain.BankTypeGovernment, IsActive: true}
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.GetByIDFunc = func(_ context.Context, id string) (*domain.Bank, error) {
		if id == "b1" {
			return bank, nil
		}
		return nil, nil
	}
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/banks/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/banks/b1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Bank
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "b1", got.ID)
	assert.Equal(t, "Bank One", got.Name)
}

func TestBankHandler_GetByID_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.GetByIDFunc = func(_ context.Context, _ string) (*domain.Bank, error) { return nil, nil }
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/banks/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/banks/none", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestBankHandler_List(t *testing.T) {
	log, _ := zap.NewDevelopment()
	list := []*domain.Bank{
		{ID: "b1", Name: "Bank A", Type: domain.BankTypePrivate, IsActive: true},
	}
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.ListFunc = func(_ context.Context, limit, offset int) ([]*domain.Bank, error) { return list, nil }
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/banks", h.List)

	req := httptest.NewRequest(http.MethodGet, "/v1/banks", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []*domain.Bank
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	require.Len(t, got, 1)
	assert.Equal(t, "b1", got[0].ID)
}

func TestBankHandler_Update(t *testing.T) {
	log, _ := zap.NewDevelopment()
	updated := &domain.Bank{ID: "b1", Name: "Bank Updated", Type: domain.BankTypePrivate, IsActive: true}
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.UpdateFunc = func(_ context.Context, id string, input domain.UpdateBankInput) (*domain.Bank, error) {
		out := *updated
		out.Name = input.Name
		out.Type = input.Type
		return &out, nil
	}
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT "+apiVersion+"/banks/{id}", h.Update)

	body := []byte(`{"name":"Bank Updated","type":"PRIVATE"}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/banks/b1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Bank
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "Bank Updated", got.Name)
}

func TestBankHandler_Update_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.UpdateFunc = func(_ context.Context, _ string, _ domain.UpdateBankInput) (*domain.Bank, error) { return nil, nil }
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT "+apiVersion+"/banks/{id}", h.Update)

	body := []byte(`{"name":"X","type":"PRIVATE"}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/banks/none", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestBankHandler_Delete(t *testing.T) {
	log, _ := zap.NewDevelopment()
	softDeleted := &domain.Bank{ID: "b1", Name: "Bank", IsActive: false}
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.DeleteFunc = func(_ context.Context, id string) (*domain.Bank, error) {
		if id == "b1" {
			return softDeleted, nil
		}
		return nil, nil
	}
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE "+apiVersion+"/banks/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/v1/banks/b1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestBankHandler_Delete_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockBankService{}
	mockSvc.DeleteFunc = func(_ context.Context, _ string) (*domain.Bank, error) { return nil, nil }
	h := NewBankHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE "+apiVersion+"/banks/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/v1/banks/none", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
