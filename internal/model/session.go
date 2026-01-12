package model

import "time"

type Session struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"userId"`
	Name      string    `json:"name" db:"name"`
	JID       string    `json:"jid,omitempty" db:"jid"`
	QRCode    string    `json:"qrCode,omitempty" db:"qrCode"`
	Connected int       `json:"connected" db:"connected"`
	Webhook   string    `json:"webhook,omitempty" db:"webhook"`
	Events    string    `json:"events,omitempty" db:"events"`
	ProxyURL  string    `json:"proxyUrl,omitempty" db:"proxyUrl"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}

type SessionCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Webhook  string `json:"webhook,omitempty"`
	Events   string `json:"events,omitempty"`
	ProxyURL string `json:"proxyUrl,omitempty"`
}

type SessionUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Webhook  *string `json:"webhook,omitempty"`
	Events   *string `json:"events,omitempty"`
	ProxyURL *string `json:"proxyUrl,omitempty"`
}



type SessionStatusResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
	LoggedIn  bool   `json:"loggedIn"`
	JID       string `json:"jid,omitempty"`
	Webhook   string `json:"webhook,omitempty"`
	Events    string `json:"events,omitempty"`
}
