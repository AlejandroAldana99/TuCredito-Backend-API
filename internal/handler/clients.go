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

// ClientHandler handles HTTP for clients.
type ClientHandler struct {
	service *service.ClientService
	log     *zap.Logger
}

// NewClientHandler returns a new ClientHandler.
func NewClientHandler(service *service.ClientService, log *zap.Logger) *ClientHandler {
	return &ClientHandler{
		service: service,
		log:     log,
	}
}

// Handler for creating a client (POST /clients).
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	// Decode the request body into a CreateClientInput
	var input domain.CreateClientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON", "INVALID_JSON", err.Error())
		return
	}

	// Check if the required fields are present
	if input.FullName == "" || input.Email == "" || input.Country == "" {
		httputil.Error(w, http.StatusBadRequest, "full_name, email, country required", "VALIDATION", "")
		return
	}

	// Create the client using the service
	client, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.log.Error("create client", zap.Error(err))
		httputil.Error(w, http.StatusInternalServerError, "failed to create client", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusCreated, client)
}

// Handler for getting a client by ID (GET /clients/:id).
func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	// Get the id from the path
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "id required", "VALIDATION", "")
		return
	}

	// Get the client by ID using the service
	client, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("get client", zap.Error(err), zap.String("id", id))
		httputil.Error(w, http.StatusInternalServerError, "failed to get client", "INTERNAL", err.Error())
		return
	}

	// Check if the client is found
	if client == nil {
		httputil.Error(w, http.StatusNotFound, "client not found", "NOT_FOUND", "")
		return
	}

	httputil.JSON(w, http.StatusOK, client)
}

// Handler for listing clients with pagination (GET /clients?limit=20&offset=0).
func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED", "")
		return
	}

	// Get the limit and offset from the query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	// List the clients using the service
	list, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		h.log.Error("list clients", zap.Error(err))
		httputil.Error(w, http.StatusInternalServerError, "failed to list clients", "INTERNAL", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, list)
}
