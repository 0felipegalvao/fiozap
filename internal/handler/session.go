package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0felipegalvao/fiozap/internal/middleware"
	"github.com/0felipegalvao/fiozap/internal/model"
	"github.com/0felipegalvao/fiozap/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

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

func (h *SessionHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		model.RespondUnauthorized(w, errors.New("user not found"))
		return
	}

	status := h.sessionService.GetStatus(user)
	model.RespondOK(w, status)
}

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
