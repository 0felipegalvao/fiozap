package model

type TextMessage struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

type ImageMessage struct {
	Phone    string `json:"phone"`
	Image    string `json:"image"`
	Caption  string `json:"caption,omitempty"`
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type AudioMessage struct {
	Phone    string `json:"phone"`
	Audio    string `json:"audio"`
	ID       string `json:"id,omitempty"`
	PTT      *bool  `json:"ptt,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type VideoMessage struct {
	Phone    string `json:"phone"`
	Video    string `json:"video"`
	Caption  string `json:"caption,omitempty"`
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type DocumentMessage struct {
	Phone    string `json:"phone"`
	Document string `json:"document"`
	FileName string `json:"filename"`
	Caption  string `json:"caption,omitempty"`
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
}

type LocationMessage struct {
	Phone     string  `json:"phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
	ID        string  `json:"id,omitempty"`
}

type ContactMessage struct {
	Phone        string `json:"phone"`
	ContactName  string `json:"contact_name"`
	ContactVCard string `json:"contact_vcard"`
	ID           string `json:"id,omitempty"`
}

type ReactionMessage struct {
	Phone     string `json:"phone"`
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

type DeleteMessage struct {
	Phone     string `json:"phone"`
	MessageID string `json:"message_id"`
}

type ConnectRequest struct {
	Subscribe []string `json:"subscribe,omitempty"`
	Immediate bool     `json:"immediate,omitempty"`
}

type WebhookRequest struct {
	WebhookURL string   `json:"webhookurl"`
	Events     []string `json:"events,omitempty"`
}

type PairPhoneRequest struct {
	Phone string `json:"phone"`
}
