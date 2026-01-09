package wameow

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"fiozap/internal/logger"
)

func getWaLogger(module string) waLog.Logger {
	return waLog.Zerolog(logger.Sub(module))
}

type Client struct {
	wac *whatsmeow.Client
}

func NewClient(ctx context.Context, dbPath string) (*Client, error) {
	dbLog := getWaLogger("database")

	container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbPath), dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlstore: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	clientLog := getWaLogger("whatsapp")
	wac := whatsmeow.NewClient(deviceStore, clientLog)

	client := &Client{wac: wac}
	wac.AddEventHandler(client.eventHandler)

	return client, nil
}

func (c *Client) Connect(ctx context.Context) error {
	if c.wac.Store.ID == nil {
		qrChan, _ := c.wac.GetQRChannel(ctx)
		if err := c.wac.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				logger.Info("Scan the QR code above to login")
			} else {
				logger.Infof("Login event: %s", evt.Event)
			}
		}
	} else {
		if err := c.wac.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		logger.Info("Connected to WhatsApp")
	}

	return nil
}

func (c *Client) Disconnect() {
	c.wac.Disconnect()
	logger.Info("Disconnected from WhatsApp")
}

func (c *Client) WaitForInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}

func (c *Client) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		logger.Infof("Received message from %s: %s", v.Info.Sender.String(), v.Message.GetConversation())
	case *events.Connected:
		logger.Info("WhatsApp connected")
	case *events.Disconnected:
		logger.Warn("WhatsApp disconnected")
	}
}

func (c *Client) GetClient() *whatsmeow.Client {
	return c.wac
}

func (c *Client) IsConnected() bool {
	return c.wac.IsConnected()
}

func (c *Client) IsLoggedIn() bool {
	return c.wac.IsLoggedIn()
}
