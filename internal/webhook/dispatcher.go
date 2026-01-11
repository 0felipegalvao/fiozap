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
	userRepo    *repository.UserRepository
	sender      *Sender
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

func NewDispatcher(webhookRepo *repository.WebhookRepository, userRepo *repository.UserRepository) *Dispatcher {
	return &Dispatcher{
		webhookRepo: webhookRepo,
		userRepo:    userRepo,
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
		user, err := d.userRepo.GetByID(event.UserID)
		if err != nil {
			logger.Warnf("User not found for webhook %d: %v", event.ID, err)
			d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		if user.Webhook == "" {
			d.webhookRepo.MarkFailed(event.ID)
			continue
		}

		if !d.shouldSendEvent(user.Events, event.EventType) {
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
		err = d.sender.Send(ctx, user.Webhook, payload)
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
	return d.webhookRepo.Create(userID, eventType, data)
}
