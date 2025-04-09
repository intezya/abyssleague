package websocket

import (
	"abysslib/jwt"
	"github.com/gorilla/websocket"
	"net/http"
	"websocket/internal/adapters/controller/http/middleware"
	"websocket/internal/adapters/controller/http/routes"
	"websocket/internal/infrastructure/hub"
)

func SetupRoute(
	mux *http.ServeMux,
	hub *hub.Hub,
	hubName string,
	jwtService jwt.Validate,
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
