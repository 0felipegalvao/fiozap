package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// Create godoc
// @Summary Create group
// @Description Create a new WhatsApp group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{name=string,participants=[]string} true "Group data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/create [post]
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
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

	result, err := h.groupService.Create(r.Context(), user.ID, session.ID, req.Name, req.Participants)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// List godoc
// @Summary List groups
// @Description List all joined WhatsApp groups
// @Tags Group
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/list [get]
func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	result, err := h.groupService.List(r.Context(), user.ID, session.ID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetInfo godoc
// @Summary Get group info
// @Description Get detailed information about a group
// @Tags Group
// @Produce json
// @Param jid query string true "Group JID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/info [get]
func (h *GroupHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	groupJID := r.URL.Query().Get("jid")
	if groupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	result, err := h.groupService.GetInfo(r.Context(), user.ID, session.ID, groupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// GetInviteLink godoc
// @Summary Get group invite link
// @Description Get the invite link for a group
// @Tags Group
// @Produce json
// @Param jid query string true "Group JID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/invitelink [get]
func (h *GroupHandler) GetInviteLink(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	groupJID := r.URL.Query().Get("jid")
	if groupJID == "" {
		model.RespondBadRequest(w, errors.New("jid is required"))
		return
	}

	result, err := h.groupService.GetInviteLink(r.Context(), user.ID, session.ID, groupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Leave godoc
// @Summary Leave group
// @Description Leave a WhatsApp group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string} true "Group JID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/leave [post]
func (h *GroupHandler) Leave(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
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

	err := h.groupService.Leave(r.Context(), user.ID, session.ID, req.GroupJID)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Left group"})
}

// UpdateParticipants godoc
// @Summary Update group participants
// @Description Add, remove, promote or demote group participants
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string,participants=[]string,action=string} true "Participants data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/updateparticipants [post]
func (h *GroupHandler) UpdateParticipants(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
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

	result, err := h.groupService.UpdateParticipants(r.Context(), user.ID, session.ID, req.GroupJID, req.Participants, req.Action)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SetName godoc
// @Summary Set group name
// @Description Update the name of a group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string,name=string} true "Group name data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/name [post]
func (h *GroupHandler) SetName(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
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

	err := h.groupService.SetName(r.Context(), user.ID, session.ID, req.GroupJID, req.Name)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group name updated"})
}

// SetTopic godoc
// @Summary Set group topic
// @Description Update the topic/description of a group
// @Tags Group
// @Accept json
// @Produce json
// @Param request body object{jid=string,topic=string} true "Group topic data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /group/topic [post]
func (h *GroupHandler) SetTopic(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
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

	err := h.groupService.SetTopic(r.Context(), user.ID, session.ID, req.GroupJID, req.Topic)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Group topic updated"})
}
