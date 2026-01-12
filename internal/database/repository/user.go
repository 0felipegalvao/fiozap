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

	maxSessions := req.MaxSessions
	if maxSessions == 0 {
		maxSessions = 5
	}

	query := `
		INSERT INTO "fzUser" ("id", "name", "token", "maxSessions")
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, id, req.Name, req.Token, maxSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.GetByID(id)
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	query := `SELECT "id", "name", "token", COALESCE("maxSessions", 5) as "maxSessions", "createdAt" FROM "fzUser" WHERE "id" = $1`

	if err := r.db.Get(&user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByToken(token string) (*model.User, error) {
	var user model.User
	query := `SELECT "id", "name", "token", COALESCE("maxSessions", 5) as "maxSessions", "createdAt" FROM "fzUser" WHERE "token" = $1`

	if err := r.db.Get(&user, query, token); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	query := `SELECT "id", "name", "token", COALESCE("maxSessions", 5) as "maxSessions", "createdAt" FROM "fzUser"`

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
	if req.MaxSessions != nil {
		user.MaxSessions = *req.MaxSessions
	}

	query := `
		UPDATE "fzUser" 
		SET "name" = $1, "token" = $2, "maxSessions" = $3
		WHERE "id" = $4
	`

	_, err = r.db.Exec(query, user.Name, user.Token, user.MaxSessions, id)
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

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
