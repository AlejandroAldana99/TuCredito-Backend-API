package cache

import "context"

// Cache provides caching, rate limiting, and distributed locking primitives.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttlSeconds int) error
	Delete(ctx context.Context, key string) error
}
