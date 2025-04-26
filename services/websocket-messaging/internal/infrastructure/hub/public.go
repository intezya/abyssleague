package hub

import (
	"context"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/entity"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/metrics"
)

func (hub *Hub) RegisterClient(client *Client) {
	hub.register <- client
}

func (hub *Hub) UnregisterClient(client *Client) {
	hub.unregister <- client
}

func (hub *Hub) GetClients(ctx context.Context) []*entity.AuthenticationData {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	result := make([]*entity.AuthenticationData, 0, len(hub.clients))

	for client, ok := range hub.clients {
		if ok {
			result = append(result, client.authentication)
		}
	}

	return result
}

func (hub *Hub) SendToUser(ctx context.Context, userId int, jsonPayload []byte) bool {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	client, exists := hub.clientsByID[userId]
	if !exists {
		return false
	}

	metrics.WebsocketMessagesSent.Inc()

	metrics.MessageSize.Observe(float64(len(jsonPayload)))

	select {
	case client.Send <- jsonPayload:
		return true
	default:
		close(client.Send)
		delete(hub.clients, client)
		delete(hub.clientsByID, userId)

		return false
	}
}

func (hub *Hub) Broadcast(ctx context.Context, message []byte) {
	hub.broadcast <- message
}
