package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/intezya/pkglib/logger"
	"net/http"
	"websocket/internal/adapters/controller/http/middleware"
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

		client := hub.NewClient(h.hub, authData, conn)
		client.Hub.RegisterClient(client)

		welcomeMsg := map[string]string{"message": "Welcome!", "user": authData.Username()}
		msgBytes, _ := json.Marshal(welcomeMsg)
		client.Send <- msgBytes

		go client.WritePump()
		go client.ReadPump()
	}
}
