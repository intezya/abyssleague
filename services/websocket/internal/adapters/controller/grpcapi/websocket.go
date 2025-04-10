package grpcapi

import (
	websocketpb "abyssproto/websocket"
	"context"
	"github.com/intezya/pkglib/itertools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"websocket/internal/infrastructure/service"
)

type WebsocketService interface {
	GetOnline(ctx context.Context) (int, error)
	GetOnlineUsers(ctx context.Context) ([]*service.OnlineUser, error)
	SendToUser(ctx context.Context, userID int, jsonPayload []byte) error
	Broadcast(ctx context.Context, jsonPayload []byte) error
}

type WebsocketHandler struct {
	websocketpb.UnimplementedWebsocketServiceServer
	websocketService WebsocketService
}

func NewWebsocketHandler(websocketService WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{websocketService: websocketService}
}

func (h *WebsocketHandler) GetOnline(
	ctx context.Context,
	_ *emptypb.Empty,
) (
	*websocketpb.GetOnlineResponse,
	error,
) {
	result, err := h.websocketService.GetOnline(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An unexpected error occured")
	}

	return &websocketpb.GetOnlineResponse{
		Online: int64(result),
	}, err
}

func (h *WebsocketHandler) GetOnlineUsers(
	ctx context.Context,
	_ *emptypb.Empty,
) (
	*websocketpb.GetOnlineUsersResponse,
	error,
) {
	result, err := h.websocketService.GetOnlineUsers(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An unexpected error occured")
	}

	users := itertools.Map(
		func(user *service.OnlineUser) *websocketpb.OnlineUser {
			return &websocketpb.OnlineUser{
				Id:         user.Id,
				Username:   user.Username,
				HardwareID: user.HardwareID,
			}
		}, result,
	)

	return &websocketpb.GetOnlineUsersResponse{
		Users: users,
	}, err
}

func (h *WebsocketHandler) SendMessage(
	ctx context.Context,
	request *websocketpb.SendMessageRequest,
) (
	*emptypb.Empty,
	error,
) {
	if request.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "UserId is required")
	}

	if request.JsonPayload == nil {
		return nil, status.Errorf(codes.InvalidArgument, "JsonPayload is required")
	}

	err := h.websocketService.SendToUser(ctx, int(request.UserId), request.JsonPayload)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not connected")
	}

	return nil, nil
}

func (h *WebsocketHandler) Broadcast(
	ctx context.Context,
	request *websocketpb.BroadcastRequest,
) (
	*emptypb.Empty,
	error,
) {
	if request.JsonPayload == nil {
		return nil, status.Errorf(codes.InvalidArgument, "JsonPayload is required")
	}

	err := h.websocketService.Broadcast(ctx, request.JsonPayload)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An unexpected error occured")
	}

	return nil, nil
}

func (h *WebsocketHandler) Setup(gRPCServer *grpc.Server) {
	websocketpb.RegisterWebsocketServiceServer(gRPCServer, h)
}
