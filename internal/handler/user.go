package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetInfo godoc
// @Summary Get user info
// @Description Get WhatsApp user info by phone numbers
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=[]string} true "Phone numbers"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /user/info [post]
func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone []string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if len(req.Phone) == 0 {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.userService.GetInfo(r.Context(), user.ID, session.ID, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// CheckUser godoc
// @Summary Check if users are on WhatsApp
// @Description Check if phone numbers are registered on WhatsApp
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=[]string} true "Phone numbers"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /user/check [post]
func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone []string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if len(req.Phone) == 0 {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.userService.CheckUser(r.Context(), user.ID, session.ID, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetAvatar godoc
// @Summary Get user avatar
// @Description Get profile picture URL for a phone number
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=string,preview=bool} true "Phone and preview option"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /user/avatar [post]
func (h *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone   string `json:"phone"`
		Preview bool   `json:"preview"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.userService.GetAvatar(r.Context(), user.ID, session.ID, req.Phone, req.Preview)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetContacts godoc
// @Summary Get contacts
// @Description Get all contacts from the connected WhatsApp account
// @Tags User
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /user/contacts [get]
func (h *UserHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.userService.GetContacts(r.Context(), user.ID, session.ID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendPresence godoc
// @Summary Send presence
// @Description Send online/offline presence status
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{presence=string} true "Presence (available/unavailable)"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /user/presence [post]
func (h *UserHandler) SendPresence(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Presence string `json:"presence"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Presence == "" {
		req.Presence = "available"
	}

	err := h.userService.SendPresence(r.Context(), user.ID, session.ID, req.Presence)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Presence sent"})
}

// ChatPresence godoc
// @Summary Send chat presence
// @Description Send typing/recording indicator to a chat
// @Tags User
// @Accept json
// @Produce json
// @Param request body object{phone=string,state=string,media=string} true "Chat presence data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/presence [post]
func (h *UserHandler) ChatPresence(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Phone string `json:"phone"`
		State string `json:"state"`
		Media string `json:"media"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	err := h.userService.ChatPresence(r.Context(), user.ID, session.ID, req.Phone, req.State, req.Media)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Chat presence sent"})
}
