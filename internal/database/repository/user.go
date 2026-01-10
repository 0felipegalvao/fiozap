package repository

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/jmoiron/sqlx"

	"fiozap/internal/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(req *model.UserCreateRequest) (*model.User, error) {
	id := generateID()

	query := `
		INSERT INTO "fzUser" ("id", "name", "token", "webhook", "events")
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, id, req.Name, req.Token, req.Webhook, req.Events)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.GetByID(id)
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	query := `SELECT "id", "name", "token", "webhook", "jid", "qrCode" as qrcode, "connected", "expiration", "events", "proxyUrl" as proxy_url FROM "fzUser" WHERE "id" = $1`

	if err := r.db.Get(&user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByToken(token string) (*model.User, error) {
	var user model.User
	query := `SELECT "id", "name", "token", "webhook", "jid", "qrCode" as qrcode, "connected", "expiration", "events", "proxyUrl" as proxy_url FROM "fzUser" WHERE "token" = $1`

	if err := r.db.Get(&user, query, token); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	query := `SELECT "id", "name", "token", "webhook", "jid", "qrCode" as qrcode, "connected", "expiration", "events", "proxyUrl" as proxy_url FROM "fzUser"`

	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Update(id string, req *model.UserUpdateRequest) (*model.User, error) {
	user, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Token != nil {
		user.Token = *req.Token
	}
	if req.Webhook != nil {
		user.Webhook = *req.Webhook
	}
	if req.Events != nil {
		user.Events = *req.Events
	}

	query := `
		UPDATE "fzUser" 
		SET "name" = $1, "token" = $2, "webhook" = $3, "events" = $4
		WHERE "id" = $5
	`

	_, err = r.db.Exec(query, user.Name, user.Token, user.Webhook, user.Events, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.GetByID(id)
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM "fzUser" WHERE "id" = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) UpdateConnected(id string, connected int) error {
	query := `UPDATE "fzUser" SET "connected" = $1 WHERE "id" = $2`
	_, err := r.db.Exec(query, connected, id)
	return err
}

func (r *UserRepository) UpdateJID(id string, jid string) error {
	query := `UPDATE "fzUser" SET "jid" = $1 WHERE "id" = $2`
	_, err := r.db.Exec(query, jid, id)
	return err
}

func (r *UserRepository) UpdateQRCode(id string, qrcode string) error {
	query := `UPDATE "fzUser" SET "qrCode" = $1 WHERE "id" = $2`
	_, err := r.db.Exec(query, qrcode, id)
	return err
}

func (r *UserRepository) UpdateWebhook(id string, webhook string, events string) error {
	query := `UPDATE "fzUser" SET "webhook" = $1, "events" = $2 WHERE "id" = $3`
	_, err := r.db.Exec(query, webhook, events, id)
	return err
}

func (r *UserRepository) GetConnectedUsers() ([]model.User, error) {
	var users []model.User
	query := `SELECT "id", "name", "token", "webhook", "jid", "qrCode" as qrcode, "connected", "expiration", "events", "proxyUrl" as proxy_url FROM "fzUser" WHERE "connected" = 1`

	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
