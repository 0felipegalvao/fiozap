package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// SendText godoc
// @Summary Send text message
// @Description Send a text message to a phone number
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.TextMessage true "Message data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/send/text [post]
func (h *MessageHandler) SendText(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.TextMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Message == "" {
		model.RespondBadRequest(w, errors.New("message is required"))
		return
	}

	result, err := h.messageService.SendText(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendImage godoc
// @Summary Send image
// @Description Send an image to a phone number
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.ImageMessage true "Image data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/send/image [post]
func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ImageMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Image == "" {
		model.RespondBadRequest(w, errors.New("image is required"))
		return
	}

	result, err := h.messageService.SendImage(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendAudio godoc
// @Summary Send audio
// @Description Send an audio file to a phone number
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.AudioMessage true "Audio data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/send/audio [post]
func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.AudioMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Audio == "" {
		model.RespondBadRequest(w, errors.New("audio is required"))
		return
	}

	result, err := h.messageService.SendAudio(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendVideo godoc
// @Summary Send video
// @Description Send a video to a phone number
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.VideoMessage true "Video data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/send/video [post]
func (h *MessageHandler) SendVideo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.VideoMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Video == "" {
		model.RespondBadRequest(w, errors.New("video is required"))
		return
	}

	result, err := h.messageService.SendVideo(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendDocument godoc
// @Summary Send document
// @Description Send a document to a phone number
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.DocumentMessage true "Document data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/send/document [post]
func (h *MessageHandler) SendDocument(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DocumentMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.Document == "" {
		model.RespondBadRequest(w, errors.New("document is required"))
		return
	}

	if req.FileName == "" {
		model.RespondBadRequest(w, errors.New("filename is required"))
		return
	}

	result, err := h.messageService.SendDocument(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendLocation godoc
// @Summary Send location
// @Description Send a location to a phone number
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.LocationMessage true "Location data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/send/location [post]
func (h *MessageHandler) SendLocation(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.LocationMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.messageService.SendLocation(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// SendContact godoc
// @Summary Send contact
// @Description Send a contact card to a phone number
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.ContactMessage true "Contact data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/send/contact [post]
func (h *MessageHandler) SendContact(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ContactMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	result, err := h.messageService.SendContact(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// React godoc
// @Summary React to message
// @Description Send a reaction to a message
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.ReactionMessage true "Reaction data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/react [post]
func (h *MessageHandler) React(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ReactionMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.MessageID == "" {
		model.RespondBadRequest(w, errors.New("message_id is required"))
		return
	}

	result, err := h.messageService.React(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Delete godoc
// @Summary Delete message
// @Description Delete a sent message
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body model.DeleteMessage true "Delete data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /chat/delete [post]
func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	session := middleware.GetSessionFromContext(r.Context())
	if user == nil || session == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.DeleteMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	if req.MessageID == "" {
		model.RespondBadRequest(w, errors.New("message_id is required"))
		return
	}

	result, err := h.messageService.Delete(r.Context(), user.ID, session.ID, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}
