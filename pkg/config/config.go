package config

import (
	"os"
	"strconv"
)

// The application configuration from environment.
type Config struct {
	HTTPPort     int
	DBConnString string
	LogLevel     string
	PProfEnabled bool
}

// Reads configuration from environment variables.
func Load() *Config {
	port, _ := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	pprof, _ := strconv.ParseBool(getEnv("PPROF_ENABLED", "true"))
	return &Config{
		HTTPPort:     port,
		DBConnString: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
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
