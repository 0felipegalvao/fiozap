package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/0felipegalvao/fiozap/internal/database/repository"
	"github.com/0felipegalvao/fiozap/internal/middleware"
	"github.com/0felipegalvao/fiozap/internal/model"
)

var supportedEventTypes = []string{
	"Message",
	"ReadReceipt",
	"HistorySync",
	"ChatPresence",
	"Presence",
	"Connected",
	"Disconnected",
	"QR",
	"LoggedOut",
	"GroupInfo",
	"JoinedGroup",
	"CallOffer",
	"All",
}

type WebhookHandler struct {
	userRepo *repository.UserRepository
}

func NewWebhookHandler(userRepo *repository.UserRepository) *WebhookHandler {
	return &WebhookHandler{userRepo: userRepo}
}

func (h *WebhookHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var events []string
	if user.Events != "" {
		events = strings.Split(user.Events, ",")
	}

	model.RespondOK(w, map[string]interface{}{
		"webhook":   user.Webhook,
		"subscribe": events,
	})
}

func (h *WebhookHandler) Set(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.WebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	var validEvents []string
	for _, event := range req.Events {
		if isValidEvent(event) {
			validEvents = append(validEvents, event)
		}
	}

	eventString := strings.Join(validEvents, ",")

	if err := h.userRepo.UpdateWebhook(user.ID, req.WebhookURL, eventString); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	middleware.InvalidateUserCache(user.Token)

	model.RespondOK(w, map[string]interface{}{
		"webhook": req.WebhookURL,
	})
}

func (h *WebhookHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		WebhookURL string   `json:"webhook"`
		Events     []string `json:"events"`
		Active     bool     `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	webhook := req.WebhookURL
	var eventString string

	if req.Active {
		var validEvents []string
		for _, event := range req.Events {
			if isValidEvent(event) {
				validEvents = append(validEvents, event)
			}
		}
		eventString = strings.Join(validEvents, ",")
	} else {
		webhook = ""
		eventString = ""
	}

	if err := h.userRepo.UpdateWebhook(user.ID, webhook, eventString); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	middleware.InvalidateUserCache(user.Token)

	model.RespondOK(w, map[string]interface{}{
		"webhook": webhook,
		"events":  strings.Split(eventString, ","),
		"active":  req.Active,
	})
}

func (h *WebhookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	if err := h.userRepo.UpdateWebhook(user.ID, "", ""); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	middleware.InvalidateUserCache(user.Token)

	model.RespondOK(w, map[string]string{
		"details": "Webhook deleted successfully",
	})
}

func isValidEvent(event string) bool {
	for _, e := range supportedEventTypes {
		if e == event {
			return true
		}
	}
	return false
}
