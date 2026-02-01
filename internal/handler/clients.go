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

type ClientHandler struct {
	service service.ClientService
	log     *zap.Logger
}

func NewClientHandler(service service.ClientService, log *zap.Logger) *ClientHandler {
	return &ClientHandler{
		service: service,
		log:     log,
	}
}

// Create creates a client (POST /clients).
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	var input domain.CreateClientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON", "INVALID_JSON", err.Error())
		return
	}

	if input.FullName == "" || input.Email == "" || input.Country == "" {
		httputil.Error(w, http.StatusBadRequest, "full_name, email, country required", "VALIDATION", "")
		return
	}

	client, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.log.Error("create client", zap.Error(err))
		httputil.Error(w, http.StatusInternalServerError, "failed to create client", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusCreated, client)
}

// GetByID gets a client by ID (GET /clients/:id).
func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "id required", "VALIDATION", "")
		return
	}

	client, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("get client", zap.Error(err), zap.String("id", id))
		httputil.Error(w, http.StatusInternalServerError, "failed to get client", "INTERNAL", err.Error())
		return
	}

	if client == nil {
		httputil.Error(w, http.StatusNotFound, "client not found", "NOT_FOUND", "")
		return
	}

	httputil.JSON(w, http.StatusOK, client)
}

// Lists clients with pagination (GET /clients).
func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
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
		h.log.Error("list clients", zap.Error(err))
		httputil.Error(w, http.StatusInternalServerError, "failed to list clients", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, list)
}
