package service

import (
	"sync"
	"websocket/internal/domain/entity"
)

type UserID = int

type sendToUser struct {
	UserID   UserID
	JsonData []byte
}

type Hub struct {
	clients    map[UserID]*entity.Client
	sendToUser chan sendToUser
	broadcast  chan []byte
	register   chan *entity.Client
	unregister chan *entity.Client
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
		sendToUser: make(chan sendToUser),
		broadcast:  make(chan []byte),
		register:   make(chan *entity.Client),
		unregister: make(chan *entity.Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.GetAuthentication().GetID()] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			clientID := client.GetAuthentication().GetID()
			if _, ok := h.clients[clientID]; ok {
				delete(h.clients, clientID)
				_ = client.CloseClient()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					_ = client.CloseClient()
					delete(h.clients, client.GetAuthentication().GetID())
				}
			}
			h.mu.Unlock()
		case message := <-h.sendToUser:
			h.mu.Lock()
			client, ok := h.clients[message.UserID]
			if ok {
				client.Send <- message.JsonData
			}
			h.mu.Unlock()
		}
	}
}
