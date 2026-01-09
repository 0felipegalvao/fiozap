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

func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
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

	result, err := h.userService.GetInfo(r.Context(), user, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
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

	result, err := h.userService.CheckUser(r.Context(), user, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
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

	result, err := h.userService.GetAvatar(r.Context(), user, req.Phone, req.Preview)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *UserHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.userService.GetContacts(r.Context(), user)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *UserHandler) SendPresence(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
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

	err := h.userService.SendPresence(r.Context(), user, req.Presence)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Presence sent"})
}

func (h *UserHandler) ChatPresence(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
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

	err := h.userService.ChatPresence(r.Context(), user, req.Phone, req.State, req.Media)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Chat presence sent"})
}
