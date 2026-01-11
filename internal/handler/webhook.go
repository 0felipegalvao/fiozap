package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"fiozap/internal/database/repository"
	"fiozap/internal/middleware"
	"fiozap/internal/model"
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

// Get godoc
// @Summary Get webhook configuration
// @Description Get current webhook URL and subscribed events
// @Tags Webhook
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /webhook [get]
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

// Set godoc
// @Summary Set webhook configuration
// @Description Configure webhook URL and events to subscribe
// @Tags Webhook
// @Accept json
// @Produce json
// @Param request body model.WebhookRequest true "Webhook configuration"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /webhook [post]
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
		"events":  validEvents,
	})
}

// Update godoc
// @Summary Update webhook configuration
// @Description Update webhook URL, events and active status
// @Tags Webhook
// @Accept json
// @Produce json
// @Param request body object{webhook=string,events=[]string,active=bool} true "Webhook update data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /webhook [put]
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

// Delete godoc
// @Summary Delete webhook configuration
// @Description Remove webhook URL and unsubscribe from all events
// @Tags Webhook
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /webhook [delete]
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
