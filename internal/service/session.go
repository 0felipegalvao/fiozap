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
	webhookRepo *repository.WebhookRepository
	clients     map[string]*wameow.Client
	mu          sync.RWMutex
	dbConnStr   string
	dispatcher  *webhook.Dispatcher
}

func NewSessionService(userRepo *repository.UserRepository, cfg *config.Config) *SessionService {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	return &SessionService{
		userRepo:  userRepo,
		clients:   make(map[string]*wameow.Client),
		dbConnStr: connStr,
	}
}

func (s *SessionService) SetWebhookRepo(repo *repository.WebhookRepository) {
	s.webhookRepo = repo
}

func (s *SessionService) SetDispatcher(d *webhook.Dispatcher) {
	s.dispatcher = d
}

func (s *SessionService) Connect(ctx context.Context, user *model.User, subscribe []string, immediate bool) (map[string]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, exists := s.clients[user.ID]; exists {
		if client.IsConnected() {
			return nil, errors.New("already connected")
		}
	}

	client, err := wameow.NewClient(ctx, s.dbConnStr, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	client.SetEventCallback(func(eventType string, data interface{}) {
		s.handleEvent(user.ID, eventType, data)
	})

	client.SetQRCallback(func(code string) {
		if err := s.userRepo.UpdateQRCode(user.ID, code); err != nil {
			logger.Warnf("Failed to update QR code: %v", err)
		}
	})

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	s.clients[user.ID] = client

	if err := s.userRepo.UpdateConnected(user.ID, 1); err != nil {
		logger.Warnf("Failed to update connected status: %v", err)
	}

	if client.IsLoggedIn() {
		jid := client.GetJID()
		if err := s.userRepo.UpdateJID(user.ID, jid.String()); err != nil {
			logger.Warnf("Failed to update JID: %v", err)
		}
	}

	return map[string]interface{}{
		"webhook": user.Webhook,
		"jid":     user.JID,
		"events":  subscribe,
		"details": "Connected!",
	}, nil
}

func (s *SessionService) handleEvent(userID, eventType string, data interface{}) {
	if s.dispatcher != nil {
		if err := s.dispatcher.Enqueue(userID, eventType, data); err != nil {
			logger.Warnf("Failed to enqueue webhook event: %v", err)
		}
	}

	if eventType == "Connected" {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if jid, ok := dataMap["jid"].(string); ok {
				if err := s.userRepo.UpdateJID(userID, jid); err != nil {
					logger.Warnf("Failed to update JID: %v", err)
				}
			}
		}
	}

	if eventType == "Disconnected" || eventType == "LoggedOut" {
		if err := s.userRepo.UpdateConnected(userID, 0); err != nil {
			logger.Warnf("Failed to update connected status: %v", err)
		}
	}
}

func (s *SessionService) Disconnect(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.clients[user.ID]
	if !exists {
		return errors.New("no session")
	}

	if !client.IsConnected() {
		return errors.New("not connected")
	}

	client.Disconnect()
	delete(s.clients, user.ID)

	if err := s.userRepo.UpdateConnected(user.ID, 0); err != nil {
		logger.Warnf("Failed to update connected status: %v", err)
	}

	return nil
}

func (s *SessionService) Logout(ctx context.Context, user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.clients[user.ID]
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

	delete(s.clients, user.ID)

	if err := s.userRepo.UpdateConnected(user.ID, 0); err != nil {
		logger.Warnf("Failed to update connected status: %v", err)
	}

	if err := s.userRepo.UpdateJID(user.ID, ""); err != nil {
		logger.Warnf("Failed to clear JID: %v", err)
	}

	return nil
}

func (s *SessionService) GetStatus(user *model.User) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	isConnected := false
	isLoggedIn := false

	if client, exists := s.clients[user.ID]; exists {
		isConnected = client.IsConnected()
		isLoggedIn = client.IsLoggedIn()
	}

	return map[string]interface{}{
		"id":        user.ID,
		"name":      user.Name,
		"connected": isConnected,
		"loggedIn":  isLoggedIn,
		"jid":       user.JID,
		"webhook":   user.Webhook,
		"events":    user.Events,
	}
}

func (s *SessionService) GetQR(user *model.User) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[user.ID]
	if !exists {
		return "", errors.New("no session, call /session/connect first")
	}

	if !client.IsConnected() {
		return "", errors.New("not connected")
	}

	if client.IsLoggedIn() {
		return "", errors.New("already logged in")
	}

	freshUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return "", err
	}

	return freshUser.QRCode, nil
}

func (s *SessionService) PairPhone(ctx context.Context, user *model.User, phone string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[user.ID]
	if !exists {
		return "", errors.New("no session, call /session/connect first")
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
	users, err := s.userRepo.GetConnectedUsers()
	if err != nil {
		logger.Errorf("Failed to get connected users: %v", err)
		return
	}

	logger.Infof("Reconnecting %d users...", len(users))

	for _, user := range users {
		go func(u model.User) {
			_, err := s.Connect(ctx, &u, nil, false)
			if err != nil {
				logger.Warnf("Failed to reconnect user %s: %v", u.ID, err)
				s.userRepo.UpdateConnected(u.ID, 0)
			} else {
				logger.Infof("User %s reconnected", u.ID)
			}
		}(user)
	}
}

func (s *SessionService) GetClient(userID string) *wameow.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[userID]
}

func (s *SessionService) GetWhatsmeowClient(userID string) *whatsmeow.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if client, exists := s.clients[userID]; exists {
		return client.GetClient()
	}
	return nil
}
