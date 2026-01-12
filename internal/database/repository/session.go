package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"fiozap/internal/model"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(userID string, req *model.SessionCreateRequest) (*model.Session, error) {
	id := generateID()

	query := `
		INSERT INTO "fzSession" ("id", "userId", "name", "webhook", "events", "proxyUrl")
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, id, userID, req.Name, req.Webhook, req.Events, req.ProxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return r.GetByID(id)
}

func (r *SessionRepository) GetByID(id string) (*model.Session, error) {
	var session model.Session
	query := `
		SELECT "id", "userId", "name", "jid", "qrCode", "connected", "webhook", "events", "proxyUrl", "createdAt"
		FROM "fzSession" 
		WHERE "id" = $1
	`

	if err := r.db.Get(&session, query, id); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetByUserAndName(userID, name string) (*model.Session, error) {
	var session model.Session
	query := `
		SELECT "id", "userId", "name", "jid", "qrCode", "connected", "webhook", "events", "proxyUrl", "createdAt"
		FROM "fzSession" 
		WHERE "userId" = $1 AND "name" = $2
	`

	if err := r.db.Get(&session, query, userID, name); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetAllByUser(userID string) ([]model.Session, error) {
	var sessions []model.Session
	query := `
		SELECT "id", "userId", "name", "jid", "qrCode", "connected", "webhook", "events", "proxyUrl", "createdAt"
		FROM "fzSession" 
		WHERE "userId" = $1
		ORDER BY "createdAt" DESC
	`

	if err := r.db.Select(&sessions, query, userID); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *SessionRepository) GetAll() ([]model.Session, error) {
	var sessions []model.Session
	query := `
		SELECT "id", "userId", "name", "jid", "qrCode", "connected", "webhook", "events", "proxyUrl", "createdAt"
		FROM "fzSession"
		ORDER BY "createdAt" DESC
	`

	if err := r.db.Select(&sessions, query); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *SessionRepository) Update(id string, req *model.SessionUpdateRequest) (*model.Session, error) {
	session, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		session.Name = *req.Name
	}
	if req.Webhook != nil {
		session.Webhook = *req.Webhook
	}
	if req.Events != nil {
		session.Events = *req.Events
	}
	if req.ProxyURL != nil {
		session.ProxyURL = *req.ProxyURL
	}

	query := `
		UPDATE "fzSession" 
		SET "name" = $1, "webhook" = $2, "events" = $3, "proxyUrl" = $4
		WHERE "id" = $5
	`

	_, err = r.db.Exec(query, session.Name, session.Webhook, session.Events, session.ProxyURL, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return r.GetByID(id)
}

func (r *SessionRepository) Delete(id string) error {
	query := `DELETE FROM "fzSession" WHERE "id" = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SessionRepository) DeleteAllByUser(userID string) error {
	query := `DELETE FROM "fzSession" WHERE "userId" = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *SessionRepository) UpdateConnected(id string, connected int) error {
	query := `UPDATE "fzSession" SET "connected" = $1 WHERE "id" = $2`
	_, err := r.db.Exec(query, connected, id)
	return err
}

func (r *SessionRepository) UpdateJID(id string, jid string) error {
	query := `UPDATE "fzSession" SET "jid" = $1 WHERE "id" = $2`
	_, err := r.db.Exec(query, jid, id)
	return err
}

func (r *SessionRepository) UpdateQRCode(id string, qrcode string) error {
	query := `UPDATE "fzSession" SET "qrCode" = $1 WHERE "id" = $2`
	_, err := r.db.Exec(query, qrcode, id)
	return err
}

func (r *SessionRepository) UpdateWebhook(id string, webhook string, events string) error {
	query := `UPDATE "fzSession" SET "webhook" = $1, "events" = $2 WHERE "id" = $3`
	_, err := r.db.Exec(query, webhook, events, id)
	return err
}

func (r *SessionRepository) GetConnectedSessions() ([]model.Session, error) {
	var sessions []model.Session
	query := `
		SELECT "id", "userId", "name", "jid", "qrCode", "connected", "webhook", "events", "proxyUrl", "createdAt"
		FROM "fzSession" 
		WHERE "connected" = 1
	`

	if err := r.db.Select(&sessions, query); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *SessionRepository) CountByUser(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM "fzSession" WHERE "userId" = $1`

	if err := r.db.Get(&count, query, userID); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SessionRepository) BelongsToUser(sessionID, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM "fzSession" WHERE "id" = $1 AND "userId" = $2`

	if err := r.db.Get(&count, query, sessionID, userID); err != nil {
		return false, err
	}

	return count > 0, nil
}
