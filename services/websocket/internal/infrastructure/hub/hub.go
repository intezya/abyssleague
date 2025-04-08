package hub

import (
	"sync"
	"time"
	"websocket/internal/domain/message"
	"websocket/internal/infrastructure/metrics"
)

const (
	writeWaitTimeout     = 5 * time.Second
	connectionTimeout    = 10 * time.Second
	connectionPingPeriod = (connectionTimeout * 9) / 10
	maxMessageSize       = 1024
)

type UserID = int

type SendToUser struct {
	UserID   UserID
	JsonData []byte
}

type Hub struct {
	clients     map[*Client]bool
	clientsByID map[int]*Client
	mu          sync.Mutex
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
	hubType     string
}

func NewHub(hubType string) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		clientsByID: make(map[int]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		hubType:     hubType,
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if existingClient, exists := h.clientsByID[client.authentication.GetID()]; exists {
				existingClient.Send <- message.DisconnectByOtherClient
				close(existingClient.Send)
				delete(h.clients, existingClient)
			}

			client.connectTime = time.Now()

			switch h.hubType {
			case "main":
				metrics.TotalMainConnections.Inc()
				metrics.ActiveMainConnections.Inc()
			case "draft":
				metrics.TotalDraftConnections.Inc()
				metrics.ActiveDraftConnections.Inc()
			}

			h.clients[client] = true
			h.clientsByID[client.authentication.GetID()] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				duration := time.Since(client.connectTime).Seconds()
				metrics.ConnectionDuration.Observe(duration)

				switch h.hubType {
				case "main":
					metrics.ActiveMainConnections.Desc()
				case "draft":
					metrics.ActiveDraftConnections.Desc()
				}

				delete(h.clients, client)
				if h.clientsByID[client.authentication.GetID()] == client {
					delete(h.clientsByID, client.authentication.GetID())
				}
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
					if h.clientsByID[client.authentication.GetID()] == client {
						delete(h.clientsByID, client.authentication.GetID())
					}
				}
			}
			h.mu.Unlock()
		}
	}
}
