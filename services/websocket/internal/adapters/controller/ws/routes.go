package ws

import (
	"abysslib/jwt"
	"github.com/gorilla/websocket"
	"net/http"
	"websocket/internal/adapters/controller/http/middleware"
	"websocket/internal/infrastructure/hub"
)

func SetupRoute(
	mux *http.ServeMux,
	hub *hub.Hub,
	hubName string,
	jwtService jwt.Validate,
) {
	autHiddleware := middleware.NewMiddleware(jwtService)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	handler := NewHandler(autHiddleware, upgrader, hub)

	mux.HandleFunc("/websocket/"+hubName, handler.GetHandler())
}
