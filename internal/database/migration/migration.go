package migration

import (
	"github.com/jmoiron/sqlx"

	"github.com/0felipegalvao/fiozap/internal/logger"
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
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			webhook TEXT DEFAULT '',
			jid TEXT DEFAULT '',
			qrcode TEXT DEFAULT '',
			connected INTEGER DEFAULT 0,
			expiration INTEGER DEFAULT 0,
			events TEXT DEFAULT '',
			proxy_url TEXT DEFAULT ''
		)
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	logger.Debug("Users table ready")
	return nil
}

func createMessageHistoryTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS message_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			chat_jid TEXT NOT NULL,
			sender_jid TEXT NOT NULL,
			message_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			message_type TEXT NOT NULL,
			text_content TEXT,
			media_link TEXT,
			quoted_message_id TEXT,
			UNIQUE(user_id, message_id)
		)
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	indexQuery := `
		CREATE INDEX IF NOT EXISTS idx_message_history_user_chat 
		ON message_history (user_id, chat_jid, timestamp DESC)
	`
	if _, err := db.Exec(indexQuery); err != nil {
		return err
	}

	logger.Debug("Message history table ready")
	return nil
}
