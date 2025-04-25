package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/http/middleware"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/adapters/controller/http/routes"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/infrastructure/hub"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/pkg/auth"
	"net/http"
)

func SetupRoute(
	mux *http.ServeMux,
	hub *hub.Hub,
	hubName string,
	jwtService *auth.JWTHelper,
) {
	authMiddleware := middleware.NewMiddleware(jwtService)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	handler := NewHandler(authMiddleware, upgrader, hub)

	mux.HandleFunc(routes.WebsocketPathPrefix+"/"+hubName, handler.GetHandler())
}
