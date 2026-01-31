package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) Client() *redis.Client {
	return r.client
}

// Creates a Redis-backed cache
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping the Redis server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Creates a cache from an existing Redis client
func NewRedisCacheFromClient(client *redis.Client) (*RedisCache, error) {
	// Ping the Redis server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Gets the value for key
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// Sets the value for key with TTL in seconds
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

// Increments the key and returns the new value
func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Expire TTL on key.
func (r *RedisCache) Expire(ctx context.Context, key string, ttlSeconds int) error {
	return r.client.Expire(ctx, key, time.Duration(ttlSeconds)*time.Second).Err()
}

// Deletes the key
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Deserializes key into v
func (r *RedisCache) GetJSON(ctx context.Context, key string, v interface{}) error {
	s, err := r.Get(ctx, key)
	if err != nil || s == "" {
		return err
	}
	return json.Unmarshal([]byte(s), v)
}

// Serializes v and stores with TTL
func (r *RedisCache) SetJSON(ctx context.Context, key string, v interface{}, ttlSeconds int) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, string(b), ttlSeconds)
}
