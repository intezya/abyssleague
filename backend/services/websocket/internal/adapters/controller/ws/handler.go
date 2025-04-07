package ws

import (
	"github.com/fasthttp/websocket"
	"log"
	"net/http"
	"websocket/internal/domain/entity"
	"websocket/internal/domain/service"
)

type Handler struct {
	authMiddleware *Middleware
	upgrader       websocket.Upgrader
	hub            *service.Hub
}

func NewHandler(
	authMiddleware *Middleware,
	upgrader websocket.Upgrader,
	hub *service.Hub,
) *Handler {
	return &Handler{
		authMiddleware: authMiddleware,
		upgrader:       upgrader,
		hub:            hub,
	}
}

func (h *Handler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authData := h.authMiddleware.JwtAuth(w, r)

		if authData == nil {
			return
		}

		conn, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading connection:", err)
			return
		}

		client := entity.NewClient(h.hub, authData, conn)
		client.Hub.RegisterClient(client)

		go client.WritePump()
	}

}
