package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tucredito/backend-api/internal/cache"
	"github.com/tucredito/backend-api/internal/decision"
	"github.com/tucredito/backend-api/internal/event"
	"github.com/tucredito/backend-api/internal/handler"
	"github.com/tucredito/backend-api/internal/metrics"
	"github.com/tucredito/backend-api/internal/middleware"
	"github.com/tucredito/backend-api/internal/repository/postgres"
	"github.com/tucredito/backend-api/internal/service"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	creditSvc  *service.CreditService
	log        *zap.Logger
}

type Config struct {
	HTTPPort     int
	DBConnString string
	RedisAddr    string
	RedisPass    string
	RedisDB      int
	Log          *zap.Logger
}

// Builds the server and wires dependencies
func New(ctx context.Context, cfg *Config) (*Server, error) {
	// Create the database pool
	pool, err := postgres.NewPool(ctx, cfg.DBConnString)
	if err != nil {
		return nil, err
	}

	// Create repositories
	clientRepo := postgres.NewClientRepository(pool)
	bankRepo := postgres.NewBankRepository(pool)
	creditRepo := postgres.NewCreditRepository(pool)

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

	// Create the publisher and engine
	publisher := event.NewMockPublisher()
	engine := decision.NewRuleEngine()
	engine.RegisterRule(decision.PaymentRangeRule{})
	engine.RegisterRule(decision.BankTypeRule{})

	// Create the services
	clientSvc := service.NewClientService(clientRepo)
	bankSvc := service.NewBankService(bankRepo)
	creditSvc := service.NewCreditService(creditRepo, clientRepo, bankRepo, c, publisher, engine, cfg.Log)

	// Create the handlers
	clientH := handler.NewClientHandler(clientSvc, cfg.Log)
	bankH := handler.NewBankHandler(bankSvc, cfg.Log)
	creditH := handler.NewCreditHandler(creditSvc, cfg.Log)

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

	// Register the credit endpoints
	mux.HandleFunc("POST /credits", creditH.Create)
	mux.HandleFunc("GET /credits", creditH.List)
	mux.HandleFunc("GET /credits/{id}", creditH.GetByID)
	mux.HandleFunc("PATCH /credits/{id}/status", creditH.UpdateStatus)
	mux.HandleFunc("PUT /credits/{id}/status", creditH.UpdateStatus)

	// Create the middleware
	var handler http.Handler = mux
	handler = middleware.Logging(cfg.Log)(handler)
	handler = middleware.RateLimit(c, 100, 60)(handler)

	// Create the HTTP server
	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.HTTPPort),
		Handler:      metricsMiddleware(handler),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		log:        cfg.Log,
	}, nil
}

// Starts the HTTP server (blocks until error or shutdown)
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.creditSvc.Shutdown()
	return s.httpServer.Shutdown(ctx)
}

// Records request count and duration per path
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.Method + " " + r.URL.Path
		metrics.IncHTTPRequest(path)
		next.ServeHTTP(w, r)
		metrics.ObserveHTTPDuration(path, time.Since(start))
	})
}
