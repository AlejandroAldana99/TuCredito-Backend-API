package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/tucredito/backend-api/internal/handler"
	"github.com/tucredito/backend-api/internal/middleware"
	"github.com/tucredito/backend-api/internal/repository/postgres"
	"github.com/tucredito/backend-api/internal/service"
	"go.uber.org/zap"
)

// Server holds HTTP server and dependencies.
type Server struct {
	httpServer *http.Server
	log        *zap.Logger
}

// Config holds server configuration.
type Config struct {
	HTTPPort     int
	DBConnString string
	RedisAddr    string
	RedisPass    string
	RedisDB      int
	Log          *zap.Logger
}

// New builds the server and wires dependencies.
func New(ctx context.Context, cfg *Config) (*Server, error) {
	// Create the database pool
	pool, err := postgres.NewPool(ctx, cfg.DBConnString)
	if err != nil {
		return nil, err
	}

	// Create repositories
	clientRepo := postgres.NewClientRepository(pool)

	// Create the services
	clientSvc := service.NewClientService(clientRepo)

	// Create the handlers
	clientH := handler.NewClientHandler(clientSvc, cfg.Log)

	// Initialize the HTTP server
	mux := http.NewServeMux()

	// Register the client endpoints
	mux.HandleFunc("POST /clients", clientH.Create)
	mux.HandleFunc("GET /clients", clientH.List)
	mux.HandleFunc("GET /clients/{id}", clientH.GetByID)

	// Create the middleware
	var handler http.Handler = mux
	handler = middleware.Logging(cfg.Log)(handler)

	// Create the HTTP server
	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.HTTPPort),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		log:        cfg.Log,
	}, nil
}

// ListenAndServe starts the HTTP server (blocks until error or shutdown).
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
