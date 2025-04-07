package service

import "context"

type OnlineUser struct {
	Id         int64
	Username   string
	HardwareID string
}

type WebsocketService struct {
	hub *Hub
}

func NewWebsocketService(hub *Hub) *WebsocketService {
	return &WebsocketService{hub: hub}
}

func (s *WebsocketService) GetOnline(ctx context.Context) (int, error) {
	return len(s.hub.clients), nil
}

func (s *WebsocketService) GetOnlineUsers(ctx context.Context) ([]*OnlineUser, error) {
	result := make([]*OnlineUser, len(s.hub.clients))
	for idx, client := range s.hub.clients {
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
	s.hub.sendToUser <- sendToUser{
		UserID:   userID,
		JsonData: jsonPayload,
	}
	return nil
}

func (s *WebsocketService) Broadcast(ctx context.Context, jsonPayload []byte) error {
	s.hub.broadcast <- jsonPayload
	return nil
}
