package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0felipegalvao/fiozap/internal/middleware"
	"github.com/0felipegalvao/fiozap/internal/model"
	"github.com/0felipegalvao/fiozap/internal/service"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		Name         string   `json:"name"`
		Participants []string `json:"participants"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Name == "" {
		model.RespondBadRequest(w, errors.New("name is required"))
		return
	}

	result, err := h.groupService.Create(r.Context(), user, req.Name, req.Participants)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.groupService.List(r.Context(), user)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *GroupHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	groupJID := r.URL.Query().Get("jid")
	if groupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	result, err := h.groupService.GetInfo(r.Context(), user, groupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *GroupHandler) GetInviteLink(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	groupJID := r.URL.Query().Get("jid")
	if groupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	result, err := h.groupService.GetInviteLink(r.Context(), user, groupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *GroupHandler) Leave(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID string `json:"jid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.groupService.Leave(r.Context(), user, req.GroupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Left group"})
}

func (h *GroupHandler) UpdateParticipants(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID     string   `json:"jid"`
		Participants []string `json:"participants"`
		Action       string   `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	if len(req.Participants) == 0 {
		model.RespondBadRequest(w, errors.New("participants is required"))
		return
	}

	if req.Action == "" {
		req.Action = "add"
	}

	result, err := h.groupService.UpdateParticipants(r.Context(), user, req.GroupJID, req.Participants, req.Action)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

func (h *GroupHandler) SetName(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID string `json:"jid"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" || req.Name == "" {
		model.RespondBadRequest(w, errors.New("jid and name are required"))
		return
	}

	err := h.groupService.SetName(r.Context(), user, req.GroupJID, req.Name)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group name updated"})
}

func (h *GroupHandler) SetTopic(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req struct {
		GroupJID string `json:"jid"`
		Topic    string `json:"topic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.GroupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	err := h.groupService.SetTopic(r.Context(), user, req.GroupJID, req.Topic)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group topic updated"})
}
