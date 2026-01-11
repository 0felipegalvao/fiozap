package wameow

import (
	"context"
	"fmt"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"fiozap/internal/logger"
)

func getWaLogger(module string) waLog.Logger {
	return waLog.Zerolog(logger.Sub(module))
}

type EventCallback func(eventType string, data interface{})

type Client struct {
	wac           *whatsmeow.Client
	userID        string
	eventCallback EventCallback
	qrCallback    func(string)
}

func NewClient(ctx context.Context, postgresConnStr string, userID string) (*Client, error) {
	dbLog := getWaLogger("database")

	container, err := sqlstore.New(ctx, "postgres", postgresConnStr, dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlstore: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	clientLog := getWaLogger("whatsapp")
	wac := whatsmeow.NewClient(deviceStore, clientLog)

	client := &Client{
		wac:    wac,
		userID: userID,
	}
	wac.AddEventHandler(client.eventHandler)

	return client, nil
}

func (c *Client) SetEventCallback(cb EventCallback) {
	c.eventCallback = cb
}

func (c *Client) SetQRCallback(cb func(string)) {
	c.qrCallback = cb
}

func (c *Client) Connect(ctx context.Context) error {
	if c.wac.Store.ID == nil {
		qrChan, _ := c.wac.GetQRChannel(ctx)
		if err := c.wac.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, logger.Writer())
					logger.Info("Scan the QR code above to login")
					if c.qrCallback != nil {
						c.qrCallback(evt.Code)
					}
					if c.eventCallback != nil {
						c.eventCallback("QR", map[string]string{"code": evt.Code})
					}
				} else {
					logger.Infof("Login event: %s", evt.Event)
				}
			}
		}()
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

func (c *Client) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		logger.Infof("Received message from %s", v.Info.Sender.String())
		if c.eventCallback != nil {
			c.eventCallback("Message", map[string]interface{}{
				"from":         v.Info.Sender.String(),
				"chat":         v.Info.Chat.String(),
				"id":           v.Info.ID,
				"timestamp":    v.Info.Timestamp.Unix(),
				"pushName":     v.Info.PushName,
				"isGroup":      v.Info.IsGroup,
				"isFromMe":     v.Info.IsFromMe,
				"text":         v.Message.GetConversation(),
				"extendedText": getExtendedText(v),
				"messageType":  getMessageType(v),
			})
		}

	case *events.Receipt:
		if c.eventCallback != nil {
			c.eventCallback("ReadReceipt", map[string]interface{}{
				"chat":       v.Chat.String(),
				"sender":     v.Sender.String(),
				"type":       string(v.Type),
				"messageIds": v.MessageIDs,
				"timestamp":  v.Timestamp.Unix(),
			})
		}

	case *events.Presence:
		if c.eventCallback != nil {
			c.eventCallback("Presence", map[string]interface{}{
				"from":        v.From.String(),
				"unavailable": v.Unavailable,
				"lastSeen":    v.LastSeen.Unix(),
			})
		}

	case *events.ChatPresence:
		if c.eventCallback != nil {
			c.eventCallback("ChatPresence", map[string]interface{}{
				"chat":   v.Chat.String(),
				"sender": v.Sender.String(),
				"state":  string(v.State),
				"media":  string(v.Media),
			})
		}

	case *events.Connected:
		logger.Info("WhatsApp connected")
		if c.eventCallback != nil {
			c.eventCallback("Connected", map[string]interface{}{
				"jid": c.wac.Store.ID.String(),
			})
		}

	case *events.Disconnected:
		logger.Warn("WhatsApp disconnected")
		if c.eventCallback != nil {
			c.eventCallback("Disconnected", nil)
		}

	case *events.LoggedOut:
		logger.Warn("WhatsApp logged out")
		if c.eventCallback != nil {
			c.eventCallback("LoggedOut", map[string]interface{}{
				"reason": v.Reason.String(),
			})
		}

	case *events.HistorySync:
		if c.eventCallback != nil {
			c.eventCallback("HistorySync", map[string]interface{}{
				"data": v.Data,
			})
		}

	case *events.CallOffer:
		if c.eventCallback != nil {
			c.eventCallback("CallOffer", map[string]interface{}{
				"from":      v.CallCreator.String(),
				"timestamp": v.Timestamp.Unix(),
				"callId":    v.CallID,
			})
		}

	case *events.GroupInfo:
		if c.eventCallback != nil {
			c.eventCallback("GroupInfo", map[string]interface{}{
				"jid":    v.JID.String(),
				"notify": v.Notify,
			})
		}

	case *events.JoinedGroup:
		if c.eventCallback != nil {
			c.eventCallback("JoinedGroup", map[string]interface{}{
				"jid":   v.JID.String(),
				"type":  v.Type,
				"name":  v.GroupInfo.Name,
				"topic": v.GroupInfo.Topic,
			})
		}
	}
}

func getExtendedText(evt *events.Message) string {
	if evt.Message != nil && evt.Message.ExtendedTextMessage != nil {
		return evt.Message.ExtendedTextMessage.GetText()
	}
	return ""
}

func getMessageType(evt *events.Message) string {
	if evt.Message == nil {
		return "unknown"
	}
	m := evt.Message
	switch {
	case m.Conversation != nil || m.ExtendedTextMessage != nil:
		return "text"
	case m.ImageMessage != nil:
		return "image"
	case m.VideoMessage != nil:
		return "video"
	case m.AudioMessage != nil:
		return "audio"
	case m.DocumentMessage != nil:
		return "document"
	case m.StickerMessage != nil:
		return "sticker"
	case m.ContactMessage != nil:
		return "contact"
	case m.LocationMessage != nil:
		return "location"
	case m.ReactionMessage != nil:
		return "reaction"
	default:
		return "unknown"
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

func (c *Client) GetJID() types.JID {
	if c.wac.Store.ID != nil {
		return *c.wac.Store.ID
	}
	return types.JID{}
}
