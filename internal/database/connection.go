package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"fiozap/internal/config"
	"fiozap/internal/logger"
)

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	if cfg.DBType == "postgres" {
		return connectPostgres(cfg)
	}
	return connectSQLite(cfg)
}

func connectPostgres(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logger.Info("Connected to PostgreSQL")
	return db, nil
}

func connectSQLite(cfg *config.Config) (*sqlx.DB, error) {
	if err := os.MkdirAll(cfg.DBPath, 0751); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	dbPath := filepath.Join(cfg.DBPath, "fiozap.db")
	dsn := dbPath + "?_pragma=foreign_keys(1)&_busy_timeout=3000"

	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite: %w", err)
	}

	logger.Infof("Connected to SQLite: %s", dbPath)
	return db, nil
}
