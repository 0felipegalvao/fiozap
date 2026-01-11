package repository

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type WebhookEvent struct {
	ID          int64           `db:"id"`
	UserID      string          `db:"userId"`
	EventType   string          `db:"eventType"`
	Payload     json.RawMessage `db:"payload"`
	Status      string          `db:"status"`
	Attempts    int             `db:"attempts"`
	LastAttempt *time.Time      `db:"lastAttempt"`
	CreatedAt   time.Time       `db:"createdAt"`
}

type WebhookRepository struct {
	db *sqlx.DB
}

func NewWebhookRepository(db *sqlx.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) Create(userID, eventType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO "fzWebhook" ("userId", "eventType", "payload", "status", "attempts", "createdAt")
		VALUES ($1, $2, $3, 'pending', 0, NOW())
	`
	_, err = r.db.Exec(query, userID, eventType, payloadBytes)
	return err
}

func (r *WebhookRepository) GetPending(limit int) ([]WebhookEvent, error) {
	var events []WebhookEvent
	query := `
		SELECT "id", "userId", "eventType", "payload", "status", "attempts", "lastAttempt", "createdAt"
		FROM "fzWebhook"
		WHERE "status" = 'pending' AND "attempts" < 3
		ORDER BY "createdAt" ASC
		LIMIT $1
	`
	err := r.db.Select(&events, query, limit)
	return events, err
}

func (r *WebhookRepository) MarkSent(id int64) error {
	query := `UPDATE "fzWebhook" SET "status" = 'sent', "lastAttempt" = NOW() WHERE "id" = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *WebhookRepository) MarkFailed(id int64) error {
	query := `
		UPDATE "fzWebhook" 
		SET "attempts" = "attempts" + 1, "lastAttempt" = NOW(),
		    "status" = CASE WHEN "attempts" >= 2 THEN 'failed' ELSE 'pending' END
		WHERE "id" = $1
	`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *WebhookRepository) DeleteOld(olderThan time.Duration) error {
	query := `DELETE FROM "fzWebhook" WHERE "createdAt" < NOW() - $1::interval`
	_, err := r.db.Exec(query, olderThan.String())
	return err
}
