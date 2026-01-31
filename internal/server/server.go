package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tucredito/backend-api/internal/cache"
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
	bankRepo := postgres.NewBankRepository(pool)

	// Create the cache
	var c cache.Cache
	var redisClient *redis.Client
	if cfg.RedisAddr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPass,
			DB:       cfg.RedisDB,
		})
		var errCache error
		c, errCache = cache.NewRedisCacheFromClient(redisClient)
		if errCache != nil {
			cfg.Log.Warn("redis unavailable, running without cache", zap.Error(errCache))
			_ = redisClient.Close()
			c = nil
			redisClient = nil
		}
	}

	// Create the services
	clientSvc := service.NewClientService(clientRepo)
	bankSvc := service.NewBankService(bankRepo)

	// Create the handlers
	clientH := handler.NewClientHandler(clientSvc, cfg.Log)
	bankH := handler.NewBankHandler(bankSvc, cfg.Log)

	// Initialize the HTTP server
	mux := http.NewServeMux()

	// Register the client endpoints
	mux.HandleFunc("POST /clients", clientH.Create)
	mux.HandleFunc("GET /clients", clientH.List)
	mux.HandleFunc("GET /clients/{id}", clientH.GetByID)

	// Register the bank endpoints
	mux.HandleFunc("POST /banks", bankH.Create)
	mux.HandleFunc("GET /banks", bankH.List)
	mux.HandleFunc("GET /banks/{id}", bankH.GetByID)

	// Create the middleware
	var handler http.Handler = mux
	handler = middleware.Logging(cfg.Log)(handler)
	handler = middleware.RateLimit(c, 100, 60)(handler)

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
