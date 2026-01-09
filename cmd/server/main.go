package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0felipegalvao/fiozap/internal/config"
	"github.com/0felipegalvao/fiozap/internal/database"
	"github.com/0felipegalvao/fiozap/internal/database/migration"
	"github.com/0felipegalvao/fiozap/internal/logger"
	"github.com/0felipegalvao/fiozap/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Init(cfg.LogLevel, cfg.LogType == "console")
	logger.Info("Starting FioZap API...")

	db, err := database.Connect(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := migration.Run(db); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	r := router.New(cfg, db)

	srv := &http.Server{
		Addr:         cfg.Address + ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  180 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Server started on %s:%s", cfg.Address, cfg.Port)
		logger.Infof("Admin token: %s", cfg.AdminToken)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	logger.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown error: %v", err)
	}

	logger.Info("Server stopped")
}
