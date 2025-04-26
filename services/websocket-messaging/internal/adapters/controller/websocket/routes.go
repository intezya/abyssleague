package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/http/middleware"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/http/routes"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/hub"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/pkg/auth"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

func SetupRoute(
	mux *http.ServeMux,
	hub *hub.Hub,
	hubName string,
	jwtService *auth.JWTHelper,
) {
	authMiddleware := middleware.NewMiddleware(jwtService)

	upgrader := websocket.Upgrader{
		ReadBufferSize:    readBufferSize,
		WriteBufferSize:   writeBufferSize,
		CheckOrigin:       func(r *http.Request) bool { return true },
		HandshakeTimeout:  0,          // Default (no timeout)
		WriteBufferPool:   nil,        // Default
		Subprotocols:      []string{}, // No subprotocols
		Error:             nil,        // Default error handler
		EnableCompression: false,      // Compression disabled
	}

	handler := NewHandler(authMiddleware, upgrader, hub)

	mux.HandleFunc(routes.WebsocketPathPrefix+"/"+hubName, handler.GetHandler())
}
