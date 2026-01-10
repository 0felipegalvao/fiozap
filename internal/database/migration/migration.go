package migration

import (
	"context"

	"go.mau.fi/util/dbutil"

	"fiozap/internal/database/migration/upgrades"
	"fiozap/internal/logger"
)

const VersionTable = "fzVersion"

func Run(ctx context.Context, db *dbutil.Database) error {
	logger.Info("Running database migrations...")

	db.UpgradeTable = upgrades.Table

	if err := db.Upgrade(ctx); err != nil {
		return err
	}

	logger.Info("Migrations completed")
	return nil
}
