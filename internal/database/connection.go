package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mau.fi/util/dbutil"

	"fiozap/internal/config"
	"fiozap/internal/database/migration"
	"fiozap/internal/logger"
)

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logger.Infof("Connected to PostgreSQL at %s:%s", cfg.DBHost, cfg.DBPort)
	return db, nil
}

func ConnectDBUtil(cfg *config.Config) (*dbutil.Database, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := dbutil.NewWithDialect(dsn, "postgres")
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	db.Owner = "fiozap"
	db.VersionTable = migration.VersionTable

	logger.Infof("Connected to PostgreSQL at %s:%s", cfg.DBHost, cfg.DBPort)
	return db, nil
}
