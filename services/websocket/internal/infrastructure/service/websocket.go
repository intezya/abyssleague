package service

import (
	"context"
	"websocket/internal/domain/entity"
	hubpackage "websocket/internal/infrastructure/hub"
)

type OnlineUser struct {
	Id         int64
	Username   string
	HardwareID string
}

type Hub interface {
	GetClients(ctx context.Context) []*entity.Client
	SendToUser(ctx context.Context, s hubpackage.SendToUser)
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
		authentication := client.GetAuthentication()
		result[idx] = &OnlineUser{
			Id:         int64(authentication.GetID()),
			Username:   authentication.GetUsername(),
			HardwareID: authentication.GetHardwareID(),
		}
	}
	return result, nil
}

func (s *WebsocketService) SendToUser(ctx context.Context, userID int, jsonPayload []byte) error {
	s.hub.SendToUser(
		ctx,
		hubpackage.SendToUser{
			UserID:   userID,
			JsonData: jsonPayload,
		},
	)
	return nil
}

func (s *WebsocketService) Broadcast(ctx context.Context, jsonPayload []byte) error {
	s.hub.Broadcast(ctx, jsonPayload)
	return nil
}
