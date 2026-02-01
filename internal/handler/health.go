package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tucredito/backend-api/pkg/httputil"
)

type HealthHandler struct {
	pool   *pgxpool.Pool
	redis  *redis.Client
	checks map[string]func(context.Context) error
	mu     sync.RWMutex
}

// Creates a new HealthHandler
func NewHealthHandler(pool *pgxpool.Pool, redisClient *redis.Client) *HealthHandler {
	h := &HealthHandler{pool: pool, redis: redisClient, checks: make(map[string]func(context.Context) error)}

	if pool != nil {
		h.checks["postgres"] = func(ctx context.Context) error { return pool.Ping(ctx) }
	}

	if redisClient != nil {
		h.checks["redis"] = func(ctx context.Context) error { return redisClient.Ping(ctx).Err() }
	}

	return h
}

// Function to check if server responds
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Function to check if all dependencies respond
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	checks := make(map[string]func(context.Context) error, len(h.checks))

	for k, v := range h.checks {
		checks[k] = v
	}

	h.mu.RUnlock()
	ctx := r.Context()
	results := make(map[string]string)
	allOk := true

	for name, fn := range checks {
		if err := fn(ctx); err != nil {
			results[name] = err.Error()
			allOk = false
		} else {
			results[name] = "ok"
		}
	}

	if allOk {
		httputil.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "checks": results})
	} else {
		httputil.JSON(w, http.StatusServiceUnavailable, map[string]interface{}{"status": "degraded", "checks": results})
	}
}
