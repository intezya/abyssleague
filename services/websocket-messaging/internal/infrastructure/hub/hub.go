package hub

import (
	"sync"
	"time"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/message"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/metrics"
	"github.com/intezya/pkglib/logger"
)

const (
	writeWaitTimeout     = 5 * time.Second
	connectionTimeout    = 10 * time.Second
	connectionPingPeriod = (connectionTimeout * 9) / 10 //nolint:mnd // that's ok
	maxMessageSize       = 1024
)

// UserID typealias.
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
		mu:          sync.Mutex{},
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		done:        make(chan struct{}),
	}
}

func (hub *Hub) GetName() string {
	return hub.name
}

func (hub *Hub) Run() {
	logger.Log.Infof("Hub %s started", hub.name)

	for {
		select {
		case <-hub.done:
			logger.Log.Infof("Hub %s received stop signal", hub.name)

			return

		case client := <-hub.register:
			hub.mu.Lock()
			if existingClient, exists := hub.clientsByID[client.authentication.ID()]; exists {
				existingClient.Send <- message.DisconnectByOtherClient
				close(existingClient.Send)
				delete(hub.clients, existingClient)
			}

			client.connectTime = time.Now()

			metrics.TotalConnections.Inc()
			metrics.ActiveConnections.Inc()

			hub.clients[client] = true
			hub.clientsByID[client.authentication.ID()] = client
			hub.mu.Unlock()

		case client := <-hub.unregister:
			hub.mu.Lock()
			if _, ok := hub.clients[client]; ok {
				duration := time.Since(client.connectTime).Seconds()

				metrics.ConnectionDuration.Observe(duration)
				metrics.ActiveConnections.Dec()

				delete(hub.clients, client)

				if hub.clientsByID[client.authentication.ID()] == client {
					delete(hub.clientsByID, client.authentication.ID())
				}

				close(client.Send)
			}
			hub.mu.Unlock()

		case msg := <-hub.broadcast:
			hub.mu.Lock()
			for client := range hub.clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(hub.clients, client)

					if hub.clientsByID[client.authentication.ID()] == client {
						delete(hub.clientsByID, client.authentication.ID())
					}
				}
			}
			hub.mu.Unlock()
		}
	}
}

func (hub *Hub) Stop() {
	logger.Log.Infof("Stopping hub %s...", hub.name)

	close(hub.done)

	const waitForClosingTime = 100 * time.Millisecond

	time.Sleep(waitForClosingTime)

	hub.mu.Lock()
	defer hub.mu.Unlock()

	for client := range hub.clients {
		close(client.Send)
	}

	hub.clients = make(map[*Client]bool)
	hub.clientsByID = make(map[int]*Client)

	logger.Log.Infof("Hub %s stopped successfully", hub.name)
}
