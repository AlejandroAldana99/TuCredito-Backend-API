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

type BankHandler struct {
	service service.BankService
	log     *zap.Logger
}

func NewBankHandler(service service.BankService, log *zap.Logger) *BankHandler {
	return &BankHandler{
		service: service,
		log:     log,
	}
}

// Create creates a bank (POST /banks).
func (h *BankHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	var input domain.CreateBankInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON", "INVALID_JSON", err.Error())
		return
	}

	if input.Name == "" || (input.Type != domain.BankTypePrivate && input.Type != domain.BankTypeGovernment) {
		httputil.Error(w, http.StatusBadRequest, "name and type (PRIVATE|GOVERNMENT) required", "VALIDATION", "")
		return
	}

	bank, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.log.Error("create bank", zap.Error(err))
		httputil.Error(w, http.StatusInternalServerError, "failed to create bank", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusCreated, bank)
}

// GetByID gets a bank by ID (GET /banks/:id).
func (h *BankHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "id required", "VALIDATION", "")
		return
	}

	bank, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("get bank", zap.Error(err), zap.String("id", id))
		httputil.Error(w, http.StatusInternalServerError, "failed to get bank", "INTERNAL", err.Error())
		return
	}

	if bank == nil {
		httputil.Error(w, http.StatusNotFound, "bank not found", "NOT_FOUND", "")
		return
	}

	httputil.JSON(w, http.StatusOK, bank)
}

// Lists banks with pagination (GET /banks).
func (h *BankHandler) List(w http.ResponseWriter, r *http.Request) {
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
		h.log.Error("list banks", zap.Error(err))
		httputil.Error(w, http.StatusInternalServerError, "failed to list banks", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, list)
}
