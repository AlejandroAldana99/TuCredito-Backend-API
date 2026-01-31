package config

import (
	"os"
	"strconv"
)

// The application configuration from environment.
type Config struct {
	HTTPPort     int
	DBConnString string
	RedisAddr    string
	RedisPass    string
	RedisDB      int
	LogLevel     string
	PProfEnabled bool
}

// Reads configuration from environment variables.
func Load() *Config {
	port, _ := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	pprof, _ := strconv.ParseBool(getEnv("PPROF_ENABLED", "true"))
	dbConnString := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPass := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	level := getEnv("LOG_LEVEL", "info")

	return &Config{
		HTTPPort:     port,
		DBConnString: dbConnString,
		RedisAddr:    redisAddr,
		RedisPass:    redisPass,
		RedisDB:      redisDB,
		LogLevel:     level,
		PProfEnabled: pprof,
	}
}

// Gets the environment variable or the default value.
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
