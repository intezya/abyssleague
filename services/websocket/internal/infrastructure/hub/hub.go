package hub

import (
	"github.com/intezya/abyssleague/services/websocket/internal/domain/message"
	"github.com/intezya/abyssleague/services/websocket/internal/infrastructure/metrics"
	"github.com/intezya/pkglib/logger"
	"sync"
	"time"
)

const (
	writeWaitTimeout     = 5 * time.Second
	connectionTimeout    = 10 * time.Second
	connectionPingPeriod = (connectionTimeout * 9) / 10
	maxMessageSize       = 1024
)

// UserID typealias
type UserID = int

type SendToUser struct {
	UserID   UserID
	JsonData []byte
}

type Hub struct {
	name        string
	clients     map[*Client]bool
	clientsByID map[int]*Client
	mu          sync.Mutex
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
	done        chan struct{}
}

func NewHub(name string) *Hub {
	return &Hub{
		name:        name,
		clients:     make(map[*Client]bool),
		clientsByID: make(map[int]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		done:        make(chan struct{}),
	}
}

func (h *Hub) GetName() string {
	return h.name
}

func (h *Hub) Run() {
	logger.Log.Infof("Hub %s started", h.name)

	for {
		select {
		case <-h.done:
			logger.Log.Infof("Hub %s received stop signal", h.name)
			return

		case client := <-h.register:
			h.mu.Lock()
			if existingClient, exists := h.clientsByID[client.authentication.ID()]; exists {
				existingClient.Send <- message.DisconnectByOtherClient
				close(existingClient.Send)
				delete(h.clients, existingClient)
			}

			client.connectTime = time.Now()
			metrics.TotalConnections.Inc()
			metrics.ActiveConnections.Inc()

			h.clients[client] = true
			h.clientsByID[client.authentication.ID()] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				duration := time.Since(client.connectTime).Seconds()

				metrics.ConnectionDuration.Observe(duration)
				metrics.ActiveConnections.Dec()

				delete(h.clients, client)
				if h.clientsByID[client.authentication.ID()] == client {
					delete(h.clientsByID, client.authentication.ID())
				}
				close(client.Send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.clients, client)
					if h.clientsByID[client.authentication.ID()] == client {
						delete(h.clientsByID, client.authentication.ID())
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Stop() {
	logger.Log.Infof("Stopping hub %s...", h.name)

	close(h.done)

	time.Sleep(100 * time.Millisecond)

	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client.Send)
	}

	h.clients = make(map[*Client]bool)
	h.clientsByID = make(map[int]*Client)

	logger.Log.Infof("Hub %s stopped successfully", h.name)
}
