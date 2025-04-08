package ws

import (
	"abysslib/logger"
	"github.com/gorilla/websocket"
	"net/http"
	"websocket/internal/adapters/controller/http/middleware"
	"websocket/internal/domain/entity"
	"websocket/internal/infrastructure/hub"
)

type Handler struct {
	authMiddleware *middleware.SecurityMiddleware
	upgrader       websocket.Upgrader
	hub            *hub.Hub
}

func NewHandler(
	authMiddleware *middleware.SecurityMiddleware,
	upgrader websocket.Upgrader,
	hub *hub.Hub,
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

		logger.Log.Info("Request headers:", r.Header)

		conn, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Log.Warn("Error upgrading connection:", err)
			return
		}

		client := entity.NewClient(h.hub, authData, conn)
		client.Hub.RegisterClient(client)

		go client.WritePump()
	}

}
