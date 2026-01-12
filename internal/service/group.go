package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type GroupService struct {
	sessionService *SessionService
}

func NewGroupService(sessionService *SessionService) *GroupService {
	return &GroupService{sessionService: sessionService}
}

func (s *GroupService) Create(ctx context.Context, userID, sessionID string, name string, participants []string) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	var jids []types.JID
	for _, p := range participants {
		jid, err := parseGroupJID(p)
		if err != nil {
			continue
		}
		jids = append(jids, jid)
	}

	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: jids,
	}

	group, err := client.CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return map[string]interface{}{
		"jid":  group.JID.String(),
		"name": group.Name,
	}, nil
}

func (s *GroupService) List(ctx context.Context, userID, sessionID string) ([]map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	groups, err := client.GetJoinedGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	var result []map[string]interface{}
	for _, g := range groups {
		result = append(result, map[string]interface{}{
			"jid":               g.JID.String(),
			"name":              g.Name,
			"topic":             g.Topic,
			"participant_count": len(g.Participants),
			"owner":             g.OwnerJID.String(),
			"created_at":        g.GroupCreated,
		})
	}

	return result, nil
}

func (s *GroupService) GetInfo(ctx context.Context, userID, sessionID string, groupJID string) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID: %w", err)
	}

	info, err := client.GetGroupInfo(ctx, jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

	var participants []map[string]interface{}
	for _, p := range info.Participants {
		participants = append(participants, map[string]interface{}{
			"jid":      p.JID.String(),
			"is_admin": p.IsAdmin,
			"is_super": p.IsSuperAdmin,
		})
	}

	return map[string]interface{}{
		"jid":          info.JID.String(),
		"name":         info.Name,
		"topic":        info.Topic,
		"owner":        info.OwnerJID.String(),
		"created_at":   info.GroupCreated,
		"participants": participants,
		"announce":     info.IsAnnounce,
		"locked":       info.IsLocked,
	}, nil
}

func (s *GroupService) GetInviteLink(ctx context.Context, userID, sessionID string, groupJID string) (map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID: %w", err)
	}

	link, err := client.GetGroupInviteLink(ctx, jid, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get invite link: %w", err)
	}

	return map[string]interface{}{
		"link": link,
	}, nil
}

func (s *GroupService) Leave(ctx context.Context, userID, sessionID string, groupJID string) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.LeaveGroup(ctx, jid)
}

func (s *GroupService) UpdateParticipants(ctx context.Context, userID, sessionID string, groupJID string, participants []string, action string) ([]map[string]interface{}, error) {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return nil, errors.New("no session")
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID: %w", err)
	}

	var jids []types.JID
	for _, p := range participants {
		pJID, err := parseGroupJID(p)
		if err != nil {
			continue
		}
		jids = append(jids, pJID)
	}

	var change whatsmeow.ParticipantChange
	switch action {
	case "add":
		change = whatsmeow.ParticipantChangeAdd
	case "remove":
		change = whatsmeow.ParticipantChangeRemove
	case "promote":
		change = whatsmeow.ParticipantChangePromote
	case "demote":
		change = whatsmeow.ParticipantChangeDemote
	default:
		change = whatsmeow.ParticipantChangeAdd
	}

	resp, err := client.UpdateGroupParticipants(ctx, jid, jids, change)
	if err != nil {
		return nil, fmt.Errorf("failed to update participants: %w", err)
	}

	var result []map[string]interface{}
	for _, r := range resp {
		result = append(result, map[string]interface{}{
			"jid":   r.JID.String(),
			"error": r.Error,
		})
	}

	return result, nil
}

func (s *GroupService) SetName(ctx context.Context, userID, sessionID string, groupJID, name string) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.SetGroupName(ctx, jid, name)
}

func (s *GroupService) SetTopic(ctx context.Context, userID, sessionID string, groupJID, topic string) error {
	client := s.sessionService.GetWhatsmeowClient(userID, sessionID)
	if client == nil {
		return errors.New("no session")
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.SetGroupTopic(ctx, jid, "", "", topic)
}

func parseGroupJID(phone string) (types.JID, error) {
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
