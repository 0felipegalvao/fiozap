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
)

type SessionService struct {
	userRepo   *repository.UserRepository
	clients    map[string]*wameow.Client
	mu         sync.RWMutex
	dbConnStr  string
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

func (s *SessionService) Connect(ctx context.Context, user *model.User, subscribe []string, immediate bool) (map[string]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, exists := s.clients[user.ID]; exists {
		if client.IsConnected() {
			return nil, errors.New("already connected")
		}
	}

	client, err := wameow.NewClient(ctx, s.dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	s.clients[user.ID] = client

	if err := s.userRepo.UpdateConnected(user.ID, 1); err != nil {
		logger.Warnf("Failed to update connected status: %v", err)
	}

	return map[string]interface{}{
		"webhook": user.Webhook,
		"jid":     user.JID,
		"events":  subscribe,
		"details": "Connected!",
	}, nil
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
		return "", errors.New("no session")
	}

	if !client.IsConnected() {
		return "", errors.New("not connected")
	}

	if client.IsLoggedIn() {
		return "", errors.New("already logged in")
	}

	return user.QRCode, nil
}

func (s *SessionService) PairPhone(ctx context.Context, user *model.User, phone string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[user.ID]
	if !exists {
		return "", errors.New("no session")
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
