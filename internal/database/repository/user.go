package repository

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/0felipegalvao/fiozap/internal/model"
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
		INSERT INTO users (id, name, token, webhook, events)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, id, req.Name, req.Token, req.Webhook, req.Events)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.GetByID(id)
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = ?`

	if err := r.db.Get(&user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByToken(token string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE token = ?`

	if err := r.db.Get(&user, query, token); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	query := `SELECT * FROM users`

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
		UPDATE users 
		SET name = ?, token = ?, webhook = ?, events = ?
		WHERE id = ?
	`

	_, err = r.db.Exec(query, user.Name, user.Token, user.Webhook, user.Events, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.GetByID(id)
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) UpdateConnected(id string, connected int) error {
	query := `UPDATE users SET connected = ? WHERE id = ?`
	_, err := r.db.Exec(query, connected, id)
	return err
}

func (r *UserRepository) UpdateJID(id string, jid string) error {
	query := `UPDATE users SET jid = ? WHERE id = ?`
	_, err := r.db.Exec(query, jid, id)
	return err
}

func (r *UserRepository) UpdateQRCode(id string, qrcode string) error {
	query := `UPDATE users SET qrcode = ? WHERE id = ?`
	_, err := r.db.Exec(query, qrcode, id)
	return err
}

func (r *UserRepository) UpdateWebhook(id string, webhook string, events string) error {
	query := `UPDATE users SET webhook = ?, events = ? WHERE id = ?`
	_, err := r.db.Exec(query, webhook, events, id)
	return err
}

func (r *UserRepository) GetConnectedUsers() ([]model.User, error) {
	var users []model.User
	query := `SELECT * FROM users WHERE connected = 1`

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
