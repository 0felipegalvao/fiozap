package model

type User struct {
	ID         string `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Token      string `json:"token" db:"token"`
	Webhook    string `json:"webhook" db:"webhook"`
	JID        string `json:"jid" db:"jid"`
	QRCode     string `json:"qrcode" db:"qrcode"`
	Connected  int    `json:"connected" db:"connected"`
	Events     string `json:"events" db:"events"`
	ProxyURL   string `json:"proxy_url,omitempty" db:"proxy_url"`
	Expiration int64  `json:"expiration" db:"expiration"`
}

type UserCreateRequest struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	Webhook string `json:"webhook,omitempty"`
	Events  string `json:"events,omitempty"`
}

type UserUpdateRequest struct {
	Name    *string `json:"name,omitempty"`
	Token   *string `json:"token,omitempty"`
	Webhook *string `json:"webhook,omitempty"`
	Events  *string `json:"events,omitempty"`
}
