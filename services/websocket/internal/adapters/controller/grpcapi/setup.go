package grpcapi

import (
	"google.golang.org/grpc"
)

func Setup(gRPCServer *grpc.Server, websocketService WebsocketService) {
	websocketHandler := NewWebsocketHandler(websocketService)
	websocketHandler.Setup(gRPCServer)
}
