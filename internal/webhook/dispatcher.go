package webhook

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"fiozap/internal/database/repository"
	"fiozap/internal/logger"
)

type Dispatcher struct {
	webhookRepo *repository.WebhookRepository
	sessionRepo *repository.SessionRepository
	sender      *Sender
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

func NewDispatcher(webhookRepo *repository.WebhookRepository, sessionRepo *repository.SessionRepository) *Dispatcher {
	return &Dispatcher{
		webhookRepo: webhookRepo,
		sessionRepo: sessionRepo,
		sender:      NewSender(),
		stopCh:      make(chan struct{}),
	}
}

func (d *Dispatcher) Start() {
	d.wg.Add(1)
	go d.processLoop()
	logger.Info("Webhook dispatcher started")
}

func (d *Dispatcher) Stop() {
	close(d.stopCh)
	d.wg.Wait()
	logger.Info("Webhook dispatcher stopped")
}

func (d *Dispatcher) processLoop() {
	defer d.wg.Done()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopCh:
			return
		case <-ticker.C:
			d.processPending()
		}
	}
}

func (d *Dispatcher) processPending() {
	events, err := d.webhookRepo.GetPending(50)
	if err != nil {
		logger.Errorf("Failed to get pending webhooks: %v", err)
		return
	}

	for _, event := range events {
		if event.SessionID == "" {
			logger.Warnf("Webhook %d has no session ID, skipping", event.ID)
			d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		session, err := d.sessionRepo.GetByID(event.SessionID)
		if err != nil {
			logger.Warnf("Session not found for webhook %d: %v", event.ID, err)
			d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		if session.Webhook == "" {
			d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		if !d.shouldSendEvent(session.Events, event.EventType) {
			d.webhookRepo.MarkSent(event.ID)
			continue
		}

		var data interface{}
		json.Unmarshal(event.Payload, &data)

		payload := &WebhookPayload{
			Event:     event.EventType,
			Timestamp: event.CreatedAt.Unix(),
			Data:      data,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = d.sender.Send(ctx, session.Webhook, payload)
		cancel()

		if err != nil {
			logger.Warnf("Failed to send webhook %d: %v", event.ID, err)
			d.webhookRepo.MarkFailed(event.ID)
		} else {
			logger.Debugf("Webhook %d sent successfully", event.ID)
			d.webhookRepo.MarkSent(event.ID)
		}
	}
}

func (d *Dispatcher) shouldSendEvent(subscribedEvents, eventType string) bool {
	if subscribedEvents == "" {
		return false
	}

	events := strings.Split(subscribedEvents, ",")
	for _, e := range events {
		if e == "All" || e == eventType {
			return true
		}
	}
	return false
}

func (d *Dispatcher) Enqueue(userID, eventType string, data interface{}) error {
	return d.webhookRepo.Create(userID, "", eventType, data)
}

func (d *Dispatcher) EnqueueSession(userID, sessionID, eventType string, data interface{}) error {
	return d.webhookRepo.Create(userID, sessionID, eventType, data)
}
