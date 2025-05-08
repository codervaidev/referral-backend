package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codervaidev/referral-backend/internal/config"
	"github.com/codervaidev/referral-backend/internal/logger"
	"github.com/codervaidev/referral-backend/internal/server"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("No .env file found")
	}

	cfg := config.Load()
	logger.Log.Info("Starting server",
		zap.String("port", cfg.Port),
		zap.String("env", cfg.Env),
	)

	srv := server.New(cfg)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Log.Fatal("Server error",
				zap.Error(err),
			)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown",
			zap.Error(err),
		)
	}

	logger.Log.Info("Server exiting")
}
