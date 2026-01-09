package model

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Success: code >= 200 && code < 300,
		Code:    code,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

func RespondError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Success: false,
		Code:    code,
		Error:   err.Error(),
	}

	json.NewEncoder(w).Encode(resp)
}

func RespondOK(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusOK, data)
}

func RespondCreated(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusCreated, data)
}

func RespondBadRequest(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusBadRequest, err)
}

func RespondUnauthorized(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusUnauthorized, err)
}

func RespondNotFound(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusNotFound, err)
}

func RespondInternalError(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusInternalServerError, err)
}
