package grpcapi

import (
	"abysslib/itertools"
	websocketContract "abyssproto/websocketgen"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"websocket/internal/domain/service"
)

//import (
//)

type WebsocketService interface {
	GetOnline(ctx context.Context) (int, error)
	GetOnlineUsers(ctx context.Context) ([]*service.OnlineUser, error)
	SendToUser(ctx context.Context, userID int, jsonPayload []byte) error
	Broadcast(ctx context.Context, jsonPayload []byte) error
}

type WebsocketHandler struct {
	websocketContract.UnimplementedWebsocketServiceServer
	websocketService WebsocketService
}

func NewWebsocketHandler(websocketService WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{websocketService: websocketService}
}

func (h *WebsocketHandler) GetOnline(
	ctx context.Context,
	_ *emptypb.Empty,
) (
	*websocketContract.GetOnlineResponse,
	error,
) {
	result, err := h.websocketService.GetOnline(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An unexpected error occured")
	}

	return &websocketContract.GetOnlineResponse{
		Online: int64(result),
	}, err
}

func (h *WebsocketHandler) GetOnlineUsers(
	ctx context.Context,
	_ *emptypb.Empty,
) (
	*websocketContract.GetOnlineUsersResponse,
	error,
) {
	result, err := h.websocketService.GetOnlineUsers(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An unexpected error occured")
	}

	users := itertools.Map(
		func(user *service.OnlineUser) *websocketContract.OnlineUser {
			return &websocketContract.OnlineUser{
				Id:         user.Id,
				Username:   user.Username,
				HardwareID: user.HardwareID,
			}
		}, result,
	)

	return &websocketContract.GetOnlineUsersResponse{
		Users: users,
	}, err
}

func (h *WebsocketHandler) SendMessage(
	ctx context.Context,
	request *websocketContract.SendMessageRequest,
) (
	*emptypb.Empty,
	error,
) {
	err := h.websocketService.SendToUser(ctx, int(request.UserId), request.JsonPayload)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An unexpected error occured")
	}

	return nil, nil
}

func (h *WebsocketHandler) Broadcast(
	ctx context.Context,
	request *websocketContract.BroadcastRequest,
) (
	*emptypb.Empty,
	error,
) {
	err := h.websocketService.Broadcast(ctx, request.JsonPayload)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An unexpected error occured")
	}

	return nil, nil
}

func (h *WebsocketHandler) Setup(gRPCServer *grpc.Server) {
	websocketContract.RegisterWebsocketServiceServer(gRPCServer, h)
}
