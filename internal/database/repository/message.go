package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Message struct {
	ID              int64     `db:"id"`
	UserID          string    `db:"userId"`
	ChatJID         string    `db:"chatJid"`
	SenderJID       string    `db:"senderJid"`
	MessageID       string    `db:"messageId"`
	Timestamp       time.Time `db:"timestamp"`
	MessageType     string    `db:"messageType"`
	TextContent     *string   `db:"textContent"`
	MediaLink       *string   `db:"mediaLink"`
	QuotedMessageID *string   `db:"quotedMessageId"`
}

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(msg *Message) error {
	query := `
		INSERT INTO "fzMessage" ("userId", "chatJid", "senderJid", "messageId", "timestamp", "messageType", "textContent", "mediaLink", "quotedMessageId")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT ("userId", "messageId") DO NOTHING
	`
	_, err := r.db.Exec(query, msg.UserID, msg.ChatJID, msg.SenderJID, msg.MessageID, msg.Timestamp, msg.MessageType, msg.TextContent, msg.MediaLink, msg.QuotedMessageID)
	return err
}

func (r *MessageRepository) GetByChat(userID, chatJID string, limit, offset int) ([]Message, error) {
	var messages []Message
	query := `
		SELECT "id", "userId", "chatJid", "senderJid", "messageId", "timestamp", "messageType", "textContent", "mediaLink", "quotedMessageId"
		FROM "fzMessage"
		WHERE "userId" = $1 AND "chatJid" = $2
		ORDER BY "timestamp" DESC
		LIMIT $3 OFFSET $4
	`
	err := r.db.Select(&messages, query, userID, chatJID, limit, offset)
	return messages, err
}

func (r *MessageRepository) GetByID(userID, messageID string) (*Message, error) {
	var msg Message
	query := `
		SELECT "id", "userId", "chatJid", "senderJid", "messageId", "timestamp", "messageType", "textContent", "mediaLink", "quotedMessageId"
		FROM "fzMessage"
		WHERE "userId" = $1 AND "messageId" = $2
	`
	err := r.db.Get(&msg, query, userID, messageID)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *MessageRepository) DeleteOld(olderThan time.Duration) error {
	query := `DELETE FROM "fzMessage" WHERE "timestamp" < NOW() - $1::interval`
	_, err := r.db.Exec(query, olderThan.String())
	return err
}

func (r *MessageRepository) CountByUser(userID string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM "fzMessage" WHERE "userId" = $1`
	err := r.db.Get(&count, query, userID)
	return count, err
}
