package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tucredito/backend-api/internal/server"
	"github.com/tucredito/backend-api/pkg/config"
	"github.com/tucredito/backend-api/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load the configuration
	cfg := config.Load()
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal("failed to create logger", zap.Error(err))
		panic(err)
	}

	defer log.Sync()

	// Create the server
	ctx := context.Background()
	srv, err := server.New(ctx, &server.Config{
		HTTPPort:     cfg.HTTPPort,
		DBConnString: cfg.DBConnString,
		RedisAddr:    cfg.RedisAddr,
		RedisPass:    cfg.RedisPass,
		RedisDB:      cfg.RedisDB,
		Log:          log,
	})
	if err != nil {
		log.Fatal("failed to create server", zap.Error(err))
	}

	// Start the pprof server
	if cfg.PProfEnabled {
		go func() {
			_ = http.ListenAndServe(":6060", nil)
		}()
		log.Info("pprof enabled on :6060")
	}

	// Start the server
	go func() {
		log.Info("server starting", zap.Int("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// Wait for a shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down")

	// Shutdown the server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", zap.Error(err))
	}
	log.Info("server stopped")
}
