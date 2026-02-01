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
	handlermocks "github.com/tucredito/backend-api/internal/handler/mocks"
	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/service"
	"go.uber.org/zap"
)

func TestCreditHandler_Create(t *testing.T) {
	log, _ := zap.NewDevelopment()
	created := &domain.Credit{
		ID: "cr1", ClientID: "c1", BankID: "b1",
		MinPayment: 100, MaxPayment: 500, TermMonths: 12,
		CreditType: domain.CreditTypeAuto, Status: domain.CreditStatusPending,
		CreatedAt: time.Now(), IsActive: true,
	}
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.CreateFunc = func(_ context.Context, input domain.CreateCreditInput) (*domain.Credit, error) {
		out := *created
		out.ClientID = input.ClientID
		out.BankID = input.BankID
		out.MinPayment = input.MinPayment
		out.MaxPayment = input.MaxPayment
		out.TermMonths = input.TermMonths
		out.CreditType = input.CreditType
		return &out, nil
	}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/credits", h.Create)

	body := []byte(`{"client_id":"c1","bank_id":"b1","min_payment":100,"max_payment":500,"term_months":12,"credit_type":"AUTO"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/credits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var got domain.Credit
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "cr1", got.ID)
	assert.Equal(t, domain.CreditTypeAuto, got.CreditType)
}

func TestCreditHandler_Create_ClientNotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.CreateFunc = func(_ context.Context, _ domain.CreateCreditInput) (*domain.Credit, error) {
		return nil, service.ErrClientNotFound
	}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/credits", h.Create)

	body := []byte(`{"client_id":"c1","bank_id":"b1","min_payment":100,"max_payment":500,"term_months":12,"credit_type":"AUTO"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/credits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var errBody struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&errBody))
	assert.Equal(t, "NOT_FOUND", errBody.Code)
}

func TestCreditHandler_Create_ValidationError(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockCreditService{}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/credits", h.Create)

	body := []byte(`{"client_id":"c1","bank_id":"b1","min_payment":500,"max_payment":100,"term_months":12,"credit_type":"AUTO"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/credits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreditHandler_GetByID(t *testing.T) {
	log, _ := zap.NewDevelopment()
	credit := &domain.Credit{
		ID: "cr1", ClientID: "c1", BankID: "b1",
		MinPayment: 100, MaxPayment: 500, TermMonths: 12,
		CreditType: domain.CreditTypeAuto, Status: domain.CreditStatusApproved,
		IsActive: true,
	}
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.GetByIDFunc = func(_ context.Context, id string) (*domain.Credit, error) {
		if id == "cr1" {
			return credit, nil
		}
		return nil, nil
	}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/credits/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/credits/cr1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Credit
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "cr1", got.ID)
	assert.Equal(t, domain.CreditStatusApproved, got.Status)
}

func TestCreditHandler_GetByID_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.GetByIDFunc = func(_ context.Context, _ string) (*domain.Credit, error) { return nil, nil }
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/credits/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/v1/credits/none", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreditHandler_List(t *testing.T) {
	log, _ := zap.NewDevelopment()
	list := []*domain.Credit{
		{ID: "cr1", ClientID: "c1", BankID: "b1", Status: domain.CreditStatusPending, IsActive: true},
	}
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.ListFunc = func(_ context.Context, limit, offset int) ([]*domain.Credit, error) { return list, nil }
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/credits", h.List)

	req := httptest.NewRequest(http.MethodGet, "/v1/credits", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []*domain.Credit
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	require.Len(t, got, 1)
	assert.Equal(t, "cr1", got[0].ID)
}

func TestCreditHandler_ListByClientID(t *testing.T) {
	log, _ := zap.NewDevelopment()
	list := []*domain.Credit{
		{ID: "cr1", ClientID: "c1", BankID: "b1", Status: domain.CreditStatusPending, IsActive: true},
	}
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.ListByClientIDFunc = func(_ context.Context, clientID string, limit, offset int) ([]*domain.Credit, error) {
		if clientID == "c1" {
			return list, nil
		}
		return nil, nil
	}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+apiVersion+"/clients/{id}/credits", h.ListByClientID)

	req := httptest.NewRequest(http.MethodGet, "/v1/clients/c1/credits", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []*domain.Credit
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	require.Len(t, got, 1)
	assert.Equal(t, "c1", got[0].ClientID)
}

func TestCreditHandler_Update(t *testing.T) {
	log, _ := zap.NewDevelopment()
	updated := &domain.Credit{
		ID: "cr1", ClientID: "c1", BankID: "b1",
		MinPayment: 200, MaxPayment: 600, TermMonths: 24,
		Status: domain.CreditStatusApproved, IsActive: true,
	}
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.UpdateFunc = func(_ context.Context, id string, input domain.UpdateCreditInput) (*domain.Credit, error) {
		out := *updated
		out.MinPayment = input.MinPayment
		out.MaxPayment = input.MaxPayment
		out.TermMonths = input.TermMonths
		out.Status = input.Status
		return &out, nil
	}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT "+apiVersion+"/credits/{id}", h.Update)

	body := []byte(`{"min_payment":200,"max_payment":600,"term_months":24,"status":"APPROVED"}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/credits/cr1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Credit
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, 200.0, got.MinPayment)
	assert.Equal(t, domain.CreditStatusApproved, got.Status)
}

func TestCreditHandler_Delete(t *testing.T) {
	log, _ := zap.NewDevelopment()
	softDeleted := &domain.Credit{ID: "cr1", ClientID: "c1", BankID: "b1", IsActive: false}
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.DeleteFunc = func(_ context.Context, id string) (*domain.Credit, error) {
		if id == "cr1" {
			return softDeleted, nil
		}
		return nil, nil
	}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE "+apiVersion+"/credits/{id}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/v1/credits/cr1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestCreditHandler_Reenable(t *testing.T) {
	log, _ := zap.NewDevelopment()
	reenabled := &domain.Credit{ID: "cr1", ClientID: "c1", BankID: "b1", IsActive: true}
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.ReenableFunc = func(_ context.Context, id string) (*domain.Credit, error) {
		if id == "cr1" {
			return reenabled, nil
		}
		return nil, nil
	}
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/credits/{id}/reenable", h.Reenable)

	req := httptest.NewRequest(http.MethodPost, "/v1/credits/cr1/reenable", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got domain.Credit
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Equal(t, "cr1", got.ID)
	assert.True(t, got.IsActive)
}

func TestCreditHandler_Reenable_NotFound(t *testing.T) {
	log, _ := zap.NewDevelopment()
	mockSvc := &handlermocks.MockCreditService{}
	mockSvc.ReenableFunc = func(_ context.Context, _ string) (*domain.Credit, error) { return nil, nil }
	h := NewCreditHandler(mockSvc, log)
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+apiVersion+"/credits/{id}/reenable", h.Reenable)

	req := httptest.NewRequest(http.MethodPost, "/v1/credits/none/reenable", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
