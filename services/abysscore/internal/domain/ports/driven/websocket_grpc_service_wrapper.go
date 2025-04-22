package drivenports

import (
	websocketpb "abyssproto/websocket"
	"context"
)

type WebsocketService interface {
	GetOnline(ctx context.Context) (*websocketpb.GetOnlineResponse, error)
	GetOnlineUsers(ctx context.Context) (*websocketpb.GetOnlineUsersResponse, error)
	SendMessage(ctx context.Context, request *websocketpb.SendMessageRequest) error
	Broadcast(ctx context.Context, request *websocketpb.BroadcastRequest) error
}
