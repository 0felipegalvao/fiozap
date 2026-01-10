package migration

import (
	"github.com/jmoiron/sqlx"

	"fiozap/internal/logger"
)

func Run(db *sqlx.DB) error {
	logger.Info("Running database migrations...")

	if err := createUsersTable(db); err != nil {
		return err
	}

	if err := createMessageHistoryTable(db); err != nil {
		return err
	}

	logger.Info("Migrations completed")
	return nil
}

func createUsersTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS "fzUser" (
			"id" VARCHAR(64) PRIMARY KEY,
			"name" VARCHAR(255) NOT NULL,
			"token" VARCHAR(255) NOT NULL UNIQUE,
			"webhook" TEXT DEFAULT '',
			"jid" VARCHAR(255) DEFAULT '',
			"qrCode" TEXT DEFAULT '',
			"connected" INTEGER DEFAULT 0,
			"expiration" BIGINT DEFAULT 0,
			"events" TEXT DEFAULT '',
			"proxyUrl" TEXT DEFAULT ''
		)
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	logger.Debug("fzUser table ready")
	return nil
}

func createMessageHistoryTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS "fzMessageHistory" (
			"id" SERIAL PRIMARY KEY,
			"userId" VARCHAR(64) NOT NULL,
			"chatJid" VARCHAR(255) NOT NULL,
			"senderJid" VARCHAR(255) NOT NULL,
			"messageId" VARCHAR(255) NOT NULL,
			"timestamp" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"messageType" VARCHAR(50) NOT NULL,
			"textContent" TEXT,
			"mediaLink" TEXT,
			"quotedMessageId" VARCHAR(255),
			UNIQUE("userId", "messageId")
		)
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	indexQuery := `
		CREATE INDEX IF NOT EXISTS "idxFzMessageHistoryUserChat" 
		ON "fzMessageHistory" ("userId", "chatJid", "timestamp" DESC)
	`
	if _, err := db.Exec(indexQuery); err != nil {
		return err
	}

	logger.Debug("fzMessageHistory table ready")
	return nil
}
