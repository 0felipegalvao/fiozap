package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.mau.fi/whatsmeow"

	"fiozap/internal/config"
	"fiozap/internal/database/repository"
	"fiozap/internal/logger"
	"fiozap/internal/model"
	"fiozap/internal/wameow"
	"fiozap/internal/webhook"
)

type SessionService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	webhookRepo *repository.WebhookRepository
	clients     map[string]*wameow.Client // key: "userId:sessionId"
	mu          sync.RWMutex
	dbConnStr   string
	dispatcher  *webhook.Dispatcher
}

func NewSessionService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, cfg *config.Config) *SessionService {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	return &SessionService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		clients:     make(map[string]*wameow.Client),
		dbConnStr:   connStr,
	}
}

func (s *SessionService) clientKey(userID, sessionID string) string {
	return fmt.Sprintf("%s:%s", userID, sessionID)
}

func (s *SessionService) SetWebhookRepo(repo *repository.WebhookRepository) {
	s.webhookRepo = repo
}

func (s *SessionService) SetDispatcher(d *webhook.Dispatcher) {
	s.dispatcher = d
}

// CRUD operations for sessions
func (s *SessionService) CreateSession(userID string, req *model.SessionCreateRequest) (*model.Session, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	count, err := s.sessionRepo.CountByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
	}

	maxSessions := user.MaxSessions
	if maxSessions == 0 {
		maxSessions = 5
	}

	if count >= maxSessions {
		return nil, fmt.Errorf("session limit reached (max: %d)", maxSessions)
	}

	return s.sessionRepo.Create(userID, req)
}

func (s *SessionService) GetSession(sessionID string) (*model.Session, error) {
	return s.sessionRepo.GetByID(sessionID)
}

func (s *SessionService) GetSessionsByUser(userID string) ([]model.Session, error) {
	return s.sessionRepo.GetAllByUser(userID)
}

func (s *SessionService) GetAllSessions() ([]model.Session, error) {
	return s.sessionRepo.GetAll()
}

func (s *SessionService) UpdateSession(sessionID string, req *model.SessionUpdateRequest) (*model.Session, error) {
	return s.sessionRepo.Update(sessionID, req)
}

func (s *SessionService) DeleteSession(userID, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.clientKey(userID, sessionID)
	if client, exists := s.clients[key]; exists {
		client.Disconnect()
		delete(s.clients, key)
	}

	return s.sessionRepo.Delete(sessionID)
}

func (s *SessionService) SessionBelongsToUser(sessionID, userID string) (bool, error) {
	return s.sessionRepo.BelongsToUser(sessionID, userID)
}

// Connection operations
func (s *SessionService) Connect(ctx context.Context, userID string, session *model.Session) (map[string]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.clientKey(userID, session.ID)

	if client, exists := s.clients[key]; exists {
		if client.IsConnected() {
			return nil, errors.New("already connected")
		}
	}

	client, err := wameow.NewClient(ctx, s.dbConnStr, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	client.SetEventCallback(func(eventType string, data interface{}) {
		s.handleEvent(userID, session.ID, eventType, data)
	})

	client.SetQRCallback(func(code string) {
		if err := s.sessionRepo.UpdateQRCode(session.ID, code); err != nil {
			logger.Warnf("Failed to update QR code: %v", err)
		}
	})

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	s.clients[key] = client

	if err := s.sessionRepo.UpdateConnected(session.ID, 1); err != nil {
		logger.Warnf("Failed to update connected status: %v", err)
	}

	if client.IsLoggedIn() {
		jid := client.GetJID()
		if err := s.sessionRepo.UpdateJID(session.ID, jid.String()); err != nil {
			logger.Warnf("Failed to update JID: %v", err)
		}
	}

	return map[string]interface{}{
		"name":    session.Name,
		"webhook": session.Webhook,
		"jid":     session.JID,
		"events":  session.Events,
		"details": "Connected!",
	}, nil
}

func (s *SessionService) handleEvent(userID, sessionID, eventType string, data interface{}) {
	if s.dispatcher != nil {
		if err := s.dispatcher.EnqueueSession(userID, sessionID, eventType, data); err != nil {
			logger.Warnf("Failed to enqueue webhook event: %v", err)
		}
	}

	if eventType == "Connected" {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if jid, ok := dataMap["jid"].(string); ok {
				if err := s.sessionRepo.UpdateJID(sessionID, jid); err != nil {
					logger.Warnf("Failed to update JID: %v", err)
				}
			}
		}
	}

	if eventType == "Disconnected" || eventType == "LoggedOut" {
		if err := s.sessionRepo.UpdateConnected(sessionID, 0); err != nil {
			logger.Warnf("Failed to update connected status: %v", err)
		}
	}
}

func (s *SessionService) Disconnect(userID string, session *model.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.clientKey(userID, session.ID)
	client, exists := s.clients[key]
	if !exists {
		return errors.New("no session")
	}

	if !client.IsConnected() {
		return errors.New("not connected")
	}

	client.Disconnect()
	delete(s.clients, key)

	if err := s.sessionRepo.UpdateConnected(session.ID, 0); err != nil {
		logger.Warnf("Failed to update connected status: %v", err)
	}

	return nil
}

func (s *SessionService) Logout(ctx context.Context, userID string, session *model.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.clientKey(userID, session.ID)
	client, exists := s.clients[key]
	if !exists {
		return errors.New("no session")
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return errors.New("not connected or not logged in")
	}

	waClient := client.GetClient()
	if err := waClient.Logout(ctx); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	delete(s.clients, key)

	if err := s.sessionRepo.UpdateConnected(session.ID, 0); err != nil {
		logger.Warnf("Failed to update connected status: %v", err)
	}

	if err := s.sessionRepo.UpdateJID(session.ID, ""); err != nil {
		logger.Warnf("Failed to clear JID: %v", err)
	}

	return nil
}

func (s *SessionService) GetStatus(userID string, session *model.Session) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.clientKey(userID, session.ID)
	isConnected := false
	isLoggedIn := false

	if client, exists := s.clients[key]; exists {
		isConnected = client.IsConnected()
		isLoggedIn = client.IsLoggedIn()
	}

	return map[string]interface{}{
		"id":        session.ID,
		"name":      session.Name,
		"connected": isConnected,
		"loggedIn":  isLoggedIn,
		"jid":       session.JID,
		"webhook":   session.Webhook,
		"events":    session.Events,
	}
}

func (s *SessionService) GetQR(userID string, session *model.Session) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.clientKey(userID, session.ID)
	client, exists := s.clients[key]
	if !exists {
		return "", errors.New("no session, call /sessions/{id}/connect first")
	}

	if !client.IsConnected() {
		return "", errors.New("not connected")
	}

	if client.IsLoggedIn() {
		return "", errors.New("already logged in")
	}

	freshSession, err := s.sessionRepo.GetByID(session.ID)
	if err != nil {
		return "", err
	}

	return freshSession.QRCode, nil
}

func (s *SessionService) PairPhone(ctx context.Context, userID string, session *model.Session, phone string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.clientKey(userID, session.ID)
	client, exists := s.clients[key]
	if !exists {
		return "", errors.New("no session, call /sessions/{id}/connect first")
	}

	waClient := client.GetClient()
	if waClient.IsLoggedIn() {
		return "", errors.New("already paired")
	}

	code, err := waClient.PairPhone(ctx, phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	return code, nil
}

func (s *SessionService) ReconnectAll(ctx context.Context) {
	sessions, err := s.sessionRepo.GetConnectedSessions()
	if err != nil {
		logger.Errorf("Failed to get connected sessions: %v", err)
		return
	}

	logger.Infof("Reconnecting %d sessions...", len(sessions))

	for _, session := range sessions {
		go func(sess model.Session) {
			_, err := s.Connect(ctx, sess.UserID, &sess)
			if err != nil {
				logger.Warnf("Failed to reconnect session %s: %v", sess.ID, err)
				s.sessionRepo.UpdateConnected(sess.ID, 0)
			} else {
				logger.Infof("Session %s reconnected", sess.ID)
			}
		}(session)
	}
}

func (s *SessionService) GetClient(userID, sessionID string) *wameow.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := s.clientKey(userID, sessionID)
	return s.clients[key]
}

func (s *SessionService) GetWhatsmeowClient(userID, sessionID string) *whatsmeow.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.clientKey(userID, sessionID)
	if client, exists := s.clients[key]; exists {
		return client.GetClient()
	}
	return nil
}
