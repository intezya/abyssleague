package drivenports

import (
	websocketpb "abyssproto/websocket"
	"context"
)

type WebsocketService interface {
	GetOnline(ctx context.Context) (*websocketpb.GetOnlineResponse, error)
	GetOnlineSoft(ctx context.Context) *websocketpb.GetOnlineResponse
	GetOnlineUsers(ctx context.Context) (*websocketpb.GetOnlineUsersResponse, error)
	GetOnlineUsersSoft(ctx context.Context) *websocketpb.GetOnlineUsersResponse
	SendMessage(ctx context.Context, request *websocketpb.SendMessageRequest) error
	Broadcast(ctx context.Context, request *websocketpb.BroadcastRequest) error
}
