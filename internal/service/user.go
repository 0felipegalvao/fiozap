package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"

	"fiozap/internal/model"
)

type UserService struct {
	sessionService *SessionService
}

func NewUserService(sessionService *SessionService) *UserService {
	return &UserService{sessionService: sessionService}
}

func (s *UserService) GetInfo(ctx context.Context, user *model.User, phones []string) ([]map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	var jids []types.JID
	for _, phone := range phones {
		jid, err := parseUserJID(phone)
		if err != nil {
			continue
		}
		jids = append(jids, jid)
	}

	if len(jids) == 0 {
		return nil, errors.New("no valid phone numbers")
	}

	info, err := client.GetUserInfo(ctx, jids)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	var result []map[string]interface{}
	for jid, userInfo := range info {
		result = append(result, map[string]interface{}{
			"jid":         jid.String(),
			"verified":    userInfo.VerifiedName != nil,
			"name":        userInfo.VerifiedName,
			"status":      userInfo.Status,
			"picture_id":  userInfo.PictureID,
			"devices":     userInfo.Devices,
		})
	}

	return result, nil
}

func (s *UserService) CheckUser(ctx context.Context, user *model.User, phones []string) ([]map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	resp, err := client.IsOnWhatsApp(ctx, phones)
	if err != nil {
		return nil, fmt.Errorf("failed to check users: %w", err)
	}

	var result []map[string]interface{}
	for _, r := range resp {
		result = append(result, map[string]interface{}{
			"query":        r.Query,
			"jid":          r.JID.String(),
			"is_in":        r.IsIn,
			"verified_name": r.VerifiedName,
		})
	}

	return result, nil
}

func (s *UserService) GetAvatar(ctx context.Context, user *model.User, phone string, preview bool) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	jid, err := parseUserJID(phone)
	if err != nil {
		return nil, err
	}

	params := &whatsmeow.GetProfilePictureParams{
		Preview: preview,
	}

	pic, err := client.GetProfilePictureInfo(ctx, jid, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar: %w", err)
	}

	if pic == nil {
		return map[string]interface{}{
			"url": "",
			"id":  "",
		}, nil
	}

	return map[string]interface{}{
		"url": pic.URL,
		"id":  pic.ID,
	}, nil
}

func (s *UserService) GetContacts(ctx context.Context, user *model.User) ([]map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return nil, errors.New("no session")
	}

	contacts, err := client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	var result []map[string]interface{}
	for jid, contact := range contacts {
		result = append(result, map[string]interface{}{
			"jid":        jid.String(),
			"name":       contact.FullName,
			"short_name": contact.FirstName,
			"push_name":  contact.PushName,
			"business":   contact.BusinessName,
		})
	}

	return result, nil
}

func (s *UserService) SendPresence(ctx context.Context, user *model.User, presence string) error {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return errors.New("no session")
	}

	var p types.Presence
	switch presence {
	case "available":
		p = types.PresenceAvailable
	case "unavailable":
		p = types.PresenceUnavailable
	default:
		p = types.PresenceAvailable
	}

	return client.SendPresence(ctx, p)
}

func (s *UserService) ChatPresence(ctx context.Context, user *model.User, phone, state, media string) error {
	client := s.sessionService.GetWhatsmeowClient(user.ID)
	if client == nil {
		return errors.New("no session")
	}

	jid, err := parseUserJID(phone)
	if err != nil {
		return err
	}

	var chatState types.ChatPresence
	switch state {
	case "composing":
		chatState = types.ChatPresenceComposing
	case "paused":
		chatState = types.ChatPresencePaused
	default:
		chatState = types.ChatPresenceComposing
	}

	var chatMedia types.ChatPresenceMedia
	switch media {
	case "audio":
		chatMedia = types.ChatPresenceMediaAudio
	default:
		chatMedia = types.ChatPresenceMediaText
	}

	return client.SendChatPresence(ctx, jid, chatState, chatMedia)
}

func parseUserJID(phone string) (types.JID, error) {
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
