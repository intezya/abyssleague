package hub

import (
	"context"
	"log"
	"sync"
	"websocket/internal/domain/entity"
)

type UserID = int

type SendToUser struct {
	UserID   UserID
	JsonData []byte
}

type Hub struct {
	clients    map[UserID]*entity.Client
	sendToUser chan SendToUser
	broadcast  chan []byte
	register   chan *entity.Client
	unregister chan *entity.Client
	quit       chan struct{}
	mu         sync.Mutex
}

func (h *Hub) RegisterClient(client *entity.Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *entity.Client) {
	h.unregister <- client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[UserID]*entity.Client),
		sendToUser: make(chan SendToUser),
		broadcast:  make(chan []byte),
		register:   make(chan *entity.Client),
		unregister: make(chan *entity.Client),
		quit:       make(chan struct{}),
	}
}

func (h *Hub) Run(ctx context.Context) {
	done := make(chan struct{})
	defer close(done)

	go func() {
		<-ctx.Done()
		h.Stop()
	}()

	for {
		select {
		case <-h.quit:
			h.mu.Lock()
			for id, client := range h.clients {
				if err := client.CloseClient(); err != nil {
					log.Printf("Error closing client %d: %v", id, err)
				}
				delete(h.clients, id)
			}
			h.mu.Unlock()
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.GetAuthentication().GetID()] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			clientID := client.GetAuthentication().GetID()
			if _, ok := h.clients[clientID]; ok {
				delete(h.clients, clientID)
				if err := client.CloseClient(); err != nil {
					log.Printf("Error closing client %d: %v", clientID, err)
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for id, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					if err := client.CloseClient(); err != nil {
						log.Printf("Error closing client %d: %v", id, err)
					}
					delete(h.clients, id)
				}
			}
			h.mu.Unlock()

		case message := <-h.sendToUser:
			h.mu.Lock()
			client, ok := h.clients[message.UserID]
			if ok {
				select {
				case client.Send <- message.JsonData:
				default:
					if err := client.CloseClient(); err != nil {
						log.Printf("Error closing client %d: %v", message.UserID, err)
					}
					delete(h.clients, message.UserID)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Stop() {
	close(h.quit)
}

func (h *Hub) GetClients(ctx context.Context) []*entity.Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]*entity.Client, len(h.clients))
	i := 0
	for _, client := range h.clients {
		result[i] = client
		i++
	}
	return result
}

func (h *Hub) SendToUser(ctx context.Context, s SendToUser) {
	select {
	case h.sendToUser <- s:
	case <-h.quit:
	}
}

func (h *Hub) Broadcast(ctx context.Context, jsonPayload []byte) {
	select {
	case h.broadcast <- jsonPayload:
	case <-h.quit:
	}
}
