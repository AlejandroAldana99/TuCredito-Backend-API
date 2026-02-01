package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tucredito/backend-api/internal/domain"
	"github.com/tucredito/backend-api/internal/service"
	"github.com/tucredito/backend-api/pkg/httputil"
	"go.uber.org/zap"
)

type CreditHandler struct {
	service service.CreditService
	log     *zap.Logger
}

func NewCreditHandler(service service.CreditService, log *zap.Logger) *CreditHandler {
	return &CreditHandler{
		service: service,
		log:     log,
	}
}

// Create creates a credit (POST /credits) - uses worker pool for concurrent processing
func (h *CreditHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	var input domain.CreateCreditInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON", "INVALID_JSON", err.Error())
		return
	}

	if input.ClientID == "" || input.BankID == "" || input.MaxPayment < input.MinPayment || input.TermMonths <= 0 {
		httputil.Error(w, http.StatusBadRequest, "client_id, bank_id, min_payment, max_payment, term_months required and valid", "VALIDATION", "")
		return
	}

	if input.CreditType != domain.CreditTypeAuto && input.CreditType != domain.CreditTypeMortgage && input.CreditType != domain.CreditTypeCommercial {
		httputil.Error(w, http.StatusBadRequest, "credit_type must be AUTO, MORTGAGE, or COMMERCIAL", "VALIDATION", "")
		return
	}

	credit, err := h.service.Create(r.Context(), input)
	if err != nil {
		switch {
		case err == service.ErrClientNotFound:
			httputil.Error(w, http.StatusNotFound, "client not found", "NOT_FOUND", "")
		case err == service.ErrBankNotFound:
			httputil.Error(w, http.StatusNotFound, "bank not found", "NOT_FOUND", "")
		case err == service.ErrInvalidInput:
			httputil.Error(w, http.StatusBadRequest, "invalid input", "VALIDATION", err.Error())
		default:
			h.log.Error("create credit", zap.Error(err))
			httputil.Error(w, http.StatusInternalServerError, "failed to create credit", "INTERNAL", err.Error())
		}
		return
	}

	httputil.JSON(w, http.StatusCreated, credit)
}

// Returns a credit by ID (GET /credits/:id)
func (h *CreditHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "id required", "VALIDATION", "")
		return
	}

	credit, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("get credit", zap.Error(err), zap.String("id", id))
		httputil.Error(w, http.StatusInternalServerError, "failed to get credit", "INTERNAL", err.Error())
		return
	}

	if credit == nil {
		httputil.Error(w, http.StatusNotFound, "credit not found", "NOT_FOUND", "")
		return
	}

	httputil.JSON(w, http.StatusOK, credit)
}

// Updates credit status (PATCH /credits/:id/status)
func (h *CreditHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "id required", "VALIDATION", "")
		return
	}

	var body domain.UpdateCreditStatusInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON", "INVALID_JSON", err.Error())
		return
	}

	if body.Status != domain.CreditStatusApproved && body.Status != domain.CreditStatusRejected && body.Status != domain.CreditStatusPending {
		httputil.Error(w, http.StatusBadRequest, "status must be PENDING, APPROVED, or REJECTED", "VALIDATION", "")
		return
	}

	credit, err := h.service.UpdateStatus(r.Context(), id, body.Status)
	if err != nil {
		h.log.Error("update credit status", zap.Error(err), zap.String("id", id))
		httputil.Error(w, http.StatusInternalServerError, "failed to update status", "INTERNAL", err.Error())
		return
	}

	if credit == nil {
		httputil.Error(w, http.StatusNotFound, "credit not found", "NOT_FOUND", "")
		return
	}

	httputil.JSON(w, http.StatusOK, credit)
}

// Lists credits with pagination (GET /credits).
func (h *CreditHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	list, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		h.log.Error("list credits", zap.Error(err))
		httputil.Error(w, http.StatusInternalServerError, "failed to list credits", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, list)
}

// Lists credits for a client (GET /clients/:id/credits).
func (h *CreditHandler) ListByClientID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	clientID := r.PathValue("id")
	if clientID == "" {
		httputil.Error(w, http.StatusBadRequest, "client id required", "VALIDATION", "")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	list, err := h.service.ListByClientID(r.Context(), clientID, limit, offset)
	if err != nil {
		h.log.Error("list credits by client", zap.Error(err), zap.String("client_id", clientID))
		httputil.Error(w, http.StatusInternalServerError, "failed to list credits", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, list)
}
