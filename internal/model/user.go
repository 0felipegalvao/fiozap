package model

import "time"

type User struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Token       string    `json:"token" db:"token"`
	MaxSessions int       `json:"maxSessions" db:"maxSessions"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
}

type UserCreateRequest struct {
	Name        string `json:"name" example:"Felipe"`
	Token       string `json:"token" example:"abc123xyz"`
	MaxSessions int    `json:"maxSessions,omitempty" example:"5"`
}

type UserUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Token       *string `json:"token,omitempty"`
	MaxSessions *int    `json:"maxSessions,omitempty"`
}
