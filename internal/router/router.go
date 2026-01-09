package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/0felipegalvao/fiozap/internal/config"
	"github.com/0felipegalvao/fiozap/internal/database/repository"
	"github.com/0felipegalvao/fiozap/internal/handler"
	"github.com/0felipegalvao/fiozap/internal/middleware"
	"github.com/0felipegalvao/fiozap/internal/service"
)

func New(cfg *config.Config, db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()

	userRepo := repository.NewUserRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(userRepo)
	adminMiddleware := middleware.NewAdminMiddleware(cfg.AdminToken)

	healthHandler := handler.NewHealthHandler()
	adminHandler := handler.NewAdminHandler(userRepo)

	sessionService := service.NewSessionService(userRepo)
	sessionHandler := handler.NewSessionHandler(sessionService)

	messageService := service.NewMessageService(sessionService)
	messageHandler := handler.NewMessageHandler(messageService)

	userService := service.NewUserService(sessionService)
	userHandler := handler.NewUserHandler(userService)

	groupService := service.NewGroupService(sessionService)
	groupHandler := handler.NewGroupHandler(groupService)

	webhookHandler := handler.NewWebhookHandler(userRepo)

	r.Use(middleware.Logging)

	r.HandleFunc("/health", healthHandler.GetHealth).Methods("GET")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(adminMiddleware.Authenticate)
	admin.HandleFunc("/users", adminHandler.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", adminHandler.ListUsers).Methods("GET")
	admin.HandleFunc("/users", adminHandler.AddUser).Methods("POST")
	admin.HandleFunc("/users/{id}", adminHandler.EditUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods("DELETE")

	api := r.PathPrefix("").Subrouter()
	api.Use(authMiddleware.Authenticate)

	// Session
	api.HandleFunc("/session/connect", sessionHandler.Connect).Methods("POST")
	api.HandleFunc("/session/disconnect", sessionHandler.Disconnect).Methods("POST")
	api.HandleFunc("/session/logout", sessionHandler.Logout).Methods("POST")
	api.HandleFunc("/session/status", sessionHandler.GetStatus).Methods("GET")
	api.HandleFunc("/session/qr", sessionHandler.GetQR).Methods("GET")
	api.HandleFunc("/session/pairphone", sessionHandler.PairPhone).Methods("POST")

	// Messages
	api.HandleFunc("/chat/send/text", messageHandler.SendText).Methods("POST")
	api.HandleFunc("/chat/send/image", messageHandler.SendImage).Methods("POST")
	api.HandleFunc("/chat/send/audio", messageHandler.SendAudio).Methods("POST")
	api.HandleFunc("/chat/send/video", messageHandler.SendVideo).Methods("POST")
	api.HandleFunc("/chat/send/document", messageHandler.SendDocument).Methods("POST")
	api.HandleFunc("/chat/send/location", messageHandler.SendLocation).Methods("POST")
	api.HandleFunc("/chat/send/contact", messageHandler.SendContact).Methods("POST")
	api.HandleFunc("/chat/react", messageHandler.React).Methods("POST")
	api.HandleFunc("/chat/delete", messageHandler.Delete).Methods("POST")

	// User
	api.HandleFunc("/user/info", userHandler.GetInfo).Methods("POST")
	api.HandleFunc("/user/check", userHandler.CheckUser).Methods("POST")
	api.HandleFunc("/user/avatar", userHandler.GetAvatar).Methods("POST")
	api.HandleFunc("/user/contacts", userHandler.GetContacts).Methods("GET")
	api.HandleFunc("/user/presence", userHandler.SendPresence).Methods("POST")
	api.HandleFunc("/chat/presence", userHandler.ChatPresence).Methods("POST")

	// Group
	api.HandleFunc("/group/create", groupHandler.Create).Methods("POST")
	api.HandleFunc("/group/list", groupHandler.List).Methods("GET")
	api.HandleFunc("/group/info", groupHandler.GetInfo).Methods("GET")
	api.HandleFunc("/group/invitelink", groupHandler.GetInviteLink).Methods("GET")
	api.HandleFunc("/group/leave", groupHandler.Leave).Methods("POST")
	api.HandleFunc("/group/updateparticipants", groupHandler.UpdateParticipants).Methods("POST")
	api.HandleFunc("/group/name", groupHandler.SetName).Methods("POST")
	api.HandleFunc("/group/topic", groupHandler.SetTopic).Methods("POST")

	// Webhook
	api.HandleFunc("/webhook", webhookHandler.Get).Methods("GET")
	api.HandleFunc("/webhook", webhookHandler.Set).Methods("POST")
	api.HandleFunc("/webhook", webhookHandler.Update).Methods("PUT")
	api.HandleFunc("/webhook", webhookHandler.Delete).Methods("DELETE")

	return r
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
