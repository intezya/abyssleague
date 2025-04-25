package hub

import (
	"context"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/entity"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/metrics"
)

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) GetClients(ctx context.Context) []*entity.AuthenticationData {
	h.mu.Lock()
	defer h.mu.Unlock()

	result := make([]*entity.AuthenticationData, 0, len(h.clients))

	for client, ok := range h.clients {
		if ok {
			result = append(result, client.authentication)
		}
	}

	return result
}

func (h *Hub) SendToUser(ctx context.Context, userId int, jsonPayload []byte) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, exists := h.clientsByID[userId]
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
		delete(h.clients, client)
		delete(h.clientsByID, userId)
		return false
	}
}

func (h *Hub) Broadcast(ctx context.Context, message []byte) {
	h.broadcast <- message
}
