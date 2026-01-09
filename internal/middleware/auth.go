package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"fiozap/internal/database/repository"
	"fiozap/internal/logger"
	"fiozap/internal/model"
)

type contextKey string

const UserContextKey contextKey = "user"

var userCache = cache.New(5*time.Minute, 10*time.Minute)

type AuthMiddleware struct {
	userRepo *repository.UserRepository
}

func NewAuthMiddleware(userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{userRepo: userRepo}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == "" {
			token = r.Header.Get("Authorization")
		}
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			model.RespondUnauthorized(w, errors.New("missing token"))
			return
		}

		var user *model.User

		if cached, found := userCache.Get(token); found {
			user = cached.(*model.User)
		} else {
			var err error
			user, err = m.userRepo.GetByToken(token)
			if err != nil {
				logger.Warnf("Invalid token: %s", token)
				model.RespondUnauthorized(w, errors.New("invalid token"))
				return
			}
			userCache.Set(token, user, cache.DefaultExpiration)
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *model.User {
	if user, ok := ctx.Value(UserContextKey).(*model.User); ok {
		return user
	}
	return nil
}

func InvalidateUserCache(token string) {
	userCache.Delete(token)
}

func UpdateUserCache(user *model.User) {
	userCache.Set(user.Token, user, cache.DefaultExpiration)
}
