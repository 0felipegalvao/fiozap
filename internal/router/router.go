package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"

	"fiozap/internal/config"
	"fiozap/internal/database/repository"
	"fiozap/internal/handler"
	"fiozap/internal/middleware"
	"fiozap/internal/service"
	"fiozap/internal/webhook"
)

type Router struct {
	mux            *mux.Router
	dispatcher     *webhook.Dispatcher
	sessionService *service.SessionService
}

func New(cfg *config.Config, db *sqlx.DB) *Router {
	r := mux.NewRouter()

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(userRepo)
	adminMiddleware := middleware.NewAdminMiddleware(cfg.AdminToken)
	sessionMiddleware := middleware.NewSessionMiddleware(sessionRepo)

	healthHandler := handler.NewHealthHandler()
	adminHandler := handler.NewAdminHandler(userRepo)

	sessionService := service.NewSessionService(userRepo, sessionRepo, cfg)
	sessionService.SetWebhookRepo(webhookRepo)

	dispatcher := webhook.NewDispatcher(webhookRepo, sessionRepo)
	sessionService.SetDispatcher(dispatcher)

	sessionHandler := handler.NewSessionHandler(sessionService)

	messageService := service.NewMessageService(sessionService)
	messageHandler := handler.NewMessageHandler(messageService)

	userService := service.NewUserService(sessionService)
	userHandler := handler.NewUserHandler(userService)

	groupService := service.NewGroupService(sessionService)
	groupHandler := handler.NewGroupHandler(groupService)

	webhookHandler := handler.NewWebhookHandler(sessionRepo)

	r.Use(cors)
	r.Use(middleware.Logging)

	r.HandleFunc("/health", healthHandler.GetHealth).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Admin routes
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(adminMiddleware.Authenticate)
	admin.HandleFunc("/users", adminHandler.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", adminHandler.ListUsers).Methods("GET")
	admin.HandleFunc("/users", adminHandler.AddUser).Methods("POST")
	admin.HandleFunc("/users/{id}", adminHandler.EditUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/sessions", sessionHandler.AdminListAllSessions).Methods("GET")

	// API routes (authenticated)
	api := r.PathPrefix("").Subrouter()
	api.Use(authMiddleware.Authenticate)

	// Sessions CRUD (user manages own sessions)
	api.HandleFunc("/sessions", sessionHandler.ListSessions).Methods("GET")
	api.HandleFunc("/sessions", sessionHandler.CreateSession).Methods("POST")

	// Session-specific routes (require session validation)
	sessionRoutes := api.PathPrefix("/sessions/{sessionId}").Subrouter()
	sessionRoutes.Use(sessionMiddleware.ValidateSession)

	// Session management
	sessionRoutes.HandleFunc("", sessionHandler.GetSession).Methods("GET")
	sessionRoutes.HandleFunc("", sessionHandler.UpdateSession).Methods("PUT")
	sessionRoutes.HandleFunc("", sessionHandler.DeleteSession).Methods("DELETE")

	// Session connection
	sessionRoutes.HandleFunc("/connect", sessionHandler.Connect).Methods("POST")
	sessionRoutes.HandleFunc("/disconnect", sessionHandler.Disconnect).Methods("POST")
	sessionRoutes.HandleFunc("/logout", sessionHandler.Logout).Methods("POST")
	sessionRoutes.HandleFunc("/status", sessionHandler.GetStatus).Methods("GET")
	sessionRoutes.HandleFunc("/qr", sessionHandler.GetQR).Methods("GET")
	sessionRoutes.HandleFunc("/pairphone", sessionHandler.PairPhone).Methods("POST")

	// Messages (per session)
	sessionRoutes.HandleFunc("/messages/text", messageHandler.SendText).Methods("POST")
	sessionRoutes.HandleFunc("/messages/image", messageHandler.SendImage).Methods("POST")
	sessionRoutes.HandleFunc("/messages/audio", messageHandler.SendAudio).Methods("POST")
	sessionRoutes.HandleFunc("/messages/video", messageHandler.SendVideo).Methods("POST")
	sessionRoutes.HandleFunc("/messages/document", messageHandler.SendDocument).Methods("POST")
	sessionRoutes.HandleFunc("/messages/location", messageHandler.SendLocation).Methods("POST")
	sessionRoutes.HandleFunc("/messages/contact", messageHandler.SendContact).Methods("POST")
	sessionRoutes.HandleFunc("/messages/reaction", messageHandler.React).Methods("POST")
	sessionRoutes.HandleFunc("/messages/delete", messageHandler.Delete).Methods("POST")

	// User operations (per session)
	sessionRoutes.HandleFunc("/user/info", userHandler.GetInfo).Methods("POST")
	sessionRoutes.HandleFunc("/user/check", userHandler.CheckUser).Methods("POST")
	sessionRoutes.HandleFunc("/user/avatar", userHandler.GetAvatar).Methods("POST")
	sessionRoutes.HandleFunc("/user/contacts", userHandler.GetContacts).Methods("GET")
	sessionRoutes.HandleFunc("/user/presence", userHandler.SendPresence).Methods("POST")
	sessionRoutes.HandleFunc("/chat/presence", userHandler.ChatPresence).Methods("POST")

	// Group operations (per session)
	sessionRoutes.HandleFunc("/group/create", groupHandler.Create).Methods("POST")
	sessionRoutes.HandleFunc("/group/list", groupHandler.List).Methods("GET")
	sessionRoutes.HandleFunc("/group/info", groupHandler.GetInfo).Methods("GET")
	sessionRoutes.HandleFunc("/group/invitelink", groupHandler.GetInviteLink).Methods("GET")
	sessionRoutes.HandleFunc("/group/leave", groupHandler.Leave).Methods("POST")
	sessionRoutes.HandleFunc("/group/updateparticipants", groupHandler.UpdateParticipants).Methods("POST")
	sessionRoutes.HandleFunc("/group/name", groupHandler.SetName).Methods("POST")
	sessionRoutes.HandleFunc("/group/topic", groupHandler.SetTopic).Methods("POST")

	// Webhook (per session)
	sessionRoutes.HandleFunc("/webhook", webhookHandler.Get).Methods("GET")
	sessionRoutes.HandleFunc("/webhook", webhookHandler.Set).Methods("POST")
	sessionRoutes.HandleFunc("/webhook", webhookHandler.Update).Methods("PUT")
	sessionRoutes.HandleFunc("/webhook", webhookHandler.Delete).Methods("DELETE")

	return &Router{
		mux:            r,
		dispatcher:     dispatcher,
		sessionService: sessionService,
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mux.ServeHTTP(w, r)
}

func (rt *Router) StartDispatcher() {
	rt.dispatcher.Start()
}

func (rt *Router) StopDispatcher() {
	rt.dispatcher.Stop()
}

func (rt *Router) GetSessionService() *service.SessionService {
	return rt.sessionService
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
