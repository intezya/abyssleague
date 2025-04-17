package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"websocket/internal/adapters/controller/http/middleware"
	"websocket/internal/adapters/controller/http/routes"
	"websocket/internal/infrastructure/hub"
	"websocket/internal/pkg/auth"
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
