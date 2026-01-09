package middleware

import (
	"errors"
	"net/http"
	"strings"

	"fiozap/internal/model"
)

type AdminMiddleware struct {
	adminToken string
}

func NewAdminMiddleware(adminToken string) *AdminMiddleware {
	return &AdminMiddleware{adminToken: adminToken}
}

func (m *AdminMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			model.RespondUnauthorized(w, errors.New("missing admin token"))
			return
		}

		if token != m.adminToken {
			model.RespondUnauthorized(w, errors.New("invalid admin token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
