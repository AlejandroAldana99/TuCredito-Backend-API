package middleware

import (
	"net/http"

	"github.com/tucredito/backend-api/internal/cache"
	"github.com/tucredito/backend-api/pkg/httputil"
)

const (
	rateLimitWindowSec = 60
	rateLimitMaxReq    = 100
	rateLimitKeyPrefix = "ratelimit:"
)

// RateLimit limits requests per client
func RateLimit(c cache.Cache, maxReq int, windowSec int) func(http.Handler) http.Handler {
	// Validate the max requests and window seconds
	if maxReq <= 0 {
		maxReq = rateLimitMaxReq
	}

	if windowSec <= 0 {
		windowSec = rateLimitWindowSec
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c == nil {
				next.ServeHTTP(w, r)
				return
			}

			clientID := r.RemoteAddr
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				clientID = xff
			}

			key := rateLimitKeyPrefix + clientID
			ctx := r.Context()
			n, err := c.Incr(ctx, key)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if n == 1 {
				_ = c.Expire(ctx, key, windowSec)
			}

			if n > int64(maxReq) {
				httputil.Error(w, http.StatusTooManyRequests, "rate limit exceeded", "RATE_LIMIT", "")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
