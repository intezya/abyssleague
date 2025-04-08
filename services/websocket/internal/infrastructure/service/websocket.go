package service

import (
	"context"
	"errors"
	"websocket/internal/domain/entity"
	hubpackage "websocket/internal/infrastructure/hub"
)

type OnlineUser struct {
	Id         int64
	Username   string
	HardwareID string
}

type Hub interface {
	GetClients(ctx context.Context) []entity.AuthenticationData
	SendToUser(ctx context.Context, userId int, jsonPayload []byte) bool
	Broadcast(ctx context.Context, jsonPayload []byte)
}

type WebsocketService struct {
	hub Hub
}

func NewWebsocketService(hub *hubpackage.Hub) *WebsocketService {
	return &WebsocketService{hub: hub}
}

func (s *WebsocketService) GetOnline(ctx context.Context) (int, error) {
	return len(s.hub.GetClients(ctx)), nil
}

func (s *WebsocketService) GetOnlineUsers(ctx context.Context) ([]*OnlineUser, error) {
	clients := s.hub.GetClients(ctx)
	result := make([]*OnlineUser, len(clients))
	for idx, client := range clients {
		result[idx] = &OnlineUser{
			Id:         int64(client.GetID()),
			Username:   client.GetUsername(),
			HardwareID: client.GetHardwareID(),
		}
	}
	return result, nil
}

func (s *WebsocketService) SendToUser(ctx context.Context, userID int, jsonPayload []byte) error {
	if !s.hub.SendToUser(ctx, userID, jsonPayload) {
		return errors.New("failed to send message to user")
	}
	return nil
}

func (s *WebsocketService) Broadcast(ctx context.Context, jsonPayload []byte) error {
	s.hub.Broadcast(ctx, jsonPayload)
	return nil
}
