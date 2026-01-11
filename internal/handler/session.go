package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/middleware"
	"fiozap/internal/model"
	"fiozap/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

// Connect godoc
// @Summary Connect WhatsApp session
// @Description Connect and start a WhatsApp session
// @Tags Session
// @Accept json
// @Produce json
// @Param request body model.ConnectRequest false "Connection options"
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /session/connect [post]
func (h *SessionHandler) Connect(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req = model.ConnectRequest{}
	}

	result, err := h.sessionService.Connect(r.Context(), user, req.Subscribe, req.Immediate)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, result)
}

// Disconnect godoc
// @Summary Disconnect WhatsApp session
// @Description Disconnect the current WhatsApp session
// @Tags Session
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /session/disconnect [post]
func (h *SessionHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	if err := h.sessionService.Disconnect(user); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Disconnected"})
}

// Logout godoc
// @Summary Logout WhatsApp session
// @Description Logout and clear WhatsApp session data
// @Tags Session
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /session/logout [post]
func (h *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	if err := h.sessionService.Logout(r.Context(), user); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "Logged out"})
}

// GetStatus godoc
// @Summary Get session status
// @Description Get the current WhatsApp session status
// @Tags Session
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /session/status [get]
func (h *SessionHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	status := h.sessionService.GetStatus(user)
	model.RespondOK(w, status)
}

// GetQR godoc
// @Summary Get QR code
// @Description Get QR code for WhatsApp authentication
// @Tags Session
// @Produce json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Security ApiKeyAuth
// @Router /session/qr [get]
func (h *SessionHandler) GetQR(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	qr, err := h.sessionService.GetQR(user)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"qrcode": qr})
}

// PairPhone godoc
// @Summary Pair with phone number
// @Description Get pairing code for phone number authentication
// @Tags Session
// @Accept json
// @Produce json
// @Param request body model.PairPhoneRequest true "Phone number"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security ApiKeyAuth
// @Router /session/pairphone [post]
func (h *SessionHandler) PairPhone(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	var req model.PairPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Phone == "" {
		model.RespondBadRequest(w, errors.New("phone is required"))
		return
	}

	code, err := h.sessionService.PairPhone(r.Context(), user, req.Phone)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"linking_code": code})
}
