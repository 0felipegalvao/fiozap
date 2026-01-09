package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/vincent-petithory/dataurl"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"github.com/0felipegalvao/fiozap/internal/logger"
	"github.com/0felipegalvao/fiozap/internal/model"
)

type MessageService struct {
	sessionService *SessionService
}

func NewMessageService(sessionService *SessionService) *MessageService {
	return &MessageService{sessionService: sessionService}
}

func (s *MessageService) SendText(ctx context.Context, user *model.User, req *model.TextMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	msg := &waE2E.Message{
		Conversation: proto.String(req.Message),
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	logger.Infof("Message sent: %s", msgID)

	return map[string]interface{}{
		"details":   "Sent",
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendImage(ctx context.Context, user *model.User, req *model.ImageMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Image, "data:image") {
		dataURL, err := dataurl.DecodeString(req.Image)
		if err != nil {
			return nil, errors.New("invalid base64 image data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("image must be base64 encoded (data:image/...)")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(filedata)
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send image: %w", err)
	}

	return map[string]interface{}{
		"details":   "Sent",
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendAudio(ctx context.Context, user *model.User, req *model.AudioMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Audio, "data:audio") {
		dataURL, err := dataurl.DecodeString(req.Audio)
		if err != nil {
			return nil, errors.New("invalid base64 audio data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("audio must be base64 encoded (data:audio/...)")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaAudio)
	if err != nil {
		return nil, fmt.Errorf("failed to upload audio: %w", err)
	}

	ptt := true
	if req.PTT != nil {
		ptt = *req.PTT
	}

	mimeType := req.MimeType
	if mimeType == "" {
		if ptt {
			mimeType = "audio/ogg; codecs=opus"
		} else {
			mimeType = "audio/mpeg"
		}
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
			PTT:           proto.Bool(ptt),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send audio: %w", err)
	}

	return map[string]interface{}{
		"details":   "Sent",
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendVideo(ctx context.Context, user *model.User, req *model.VideoMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Video, "data:video") {
		dataURL, err := dataurl.DecodeString(req.Video)
		if err != nil {
			return nil, errors.New("invalid base64 video data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("video must be base64 encoded (data:video/...)")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaVideo)
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %w", err)
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(filedata)
	}

	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send video: %w", err)
	}

	return map[string]interface{}{
		"details":   "Sent",
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendDocument(ctx context.Context, user *model.User, req *model.DocumentMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	var filedata []byte
	if strings.HasPrefix(req.Document, "data:") {
		dataURL, err := dataurl.DecodeString(req.Document)
		if err != nil {
			return nil, errors.New("invalid base64 document data")
		}
		filedata = dataURL.Data
	} else {
		return nil, errors.New("document must be base64 encoded")
	}

	uploaded, err := client.Upload(ctx, filedata, whatsmeow.MediaDocument)
	if err != nil {
		return nil, fmt.Errorf("failed to upload document: %w", err)
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(filedata)
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(filedata))),
			FileName:      proto.String(req.FileName),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send document: %w", err)
	}

	return map[string]interface{}{
		"details":   "Sent",
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendLocation(ctx context.Context, user *model.User, req *model.LocationMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(req.Latitude),
			DegreesLongitude: proto.Float64(req.Longitude),
			Name:             proto.String(req.Name),
			Address:          proto.String(req.Address),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send location: %w", err)
	}

	return map[string]interface{}{
		"details":   "Sent",
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) SendContact(ctx context.Context, user *model.User, req *model.ContactMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GenerateMessageID()
	}

	msg := &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: proto.String(req.ContactName),
			Vcard:       proto.String(req.ContactVCard),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("failed to send contact: %w", err)
	}

	return map[string]interface{}{
		"details":   "Sent",
		"timestamp": resp.Timestamp.Unix(),
		"id":        msgID,
	}, nil
}

func (s *MessageService) React(ctx context.Context, user *model.User, req *model.ReactionMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	msg := &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(recipient.String()),
				ID:        proto.String(req.MessageID),
				FromMe:    proto.Bool(false),
			},
			Text: proto.String(req.Emoji),
		},
	}

	resp, err := client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send reaction: %w", err)
	}

	return map[string]interface{}{
		"details":   "Reacted",
		"timestamp": resp.Timestamp.Unix(),
	}, nil
}

func (s *MessageService) Delete(ctx context.Context, user *model.User, req *model.DeleteMessage) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, err
	}

	resp, err := client.RevokeMessage(ctx, recipient, types.MessageID(req.MessageID))
	if err != nil {
		return nil, fmt.Errorf("failed to delete message: %w", err)
	}

	return map[string]interface{}{
		"details":   "Deleted",
		"timestamp": resp.Timestamp.Unix(),
	}, nil
}

func parseJID(phone string) (types.JID, error) {
	if phone == "" {
		return types.JID{}, errors.New("phone is required")
	}

	if phone[0] == '+' {
		phone = phone[1:]
	}

	if !strings.Contains(phone, "@") {
		return types.NewJID(phone, types.DefaultUserServer), nil
	}

	jid, err := types.ParseJID(phone)
	if err != nil {
		return types.JID{}, fmt.Errorf("invalid JID: %w", err)
	}

	return jid, nil
}
