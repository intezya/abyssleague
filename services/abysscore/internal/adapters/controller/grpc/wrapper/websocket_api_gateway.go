package wrapper

import (
	"abysscore/internal/adapters/controller/grpc/factory"
	"context"
	"time"

	websocketpb "abyssproto/websocket"
	"github.com/intezya/pkglib/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WebsocketServiceWrapper struct {
	factory     *factory.GrpcClientFactory
	serviceAddr string
	client      websocketpb.WebsocketServiceClient
	timeout     time.Duration
}

func NewWebsocketServiceWrapper(factory *factory.GrpcClientFactory, serviceAddr string) *WebsocketServiceWrapper {
	return &WebsocketServiceWrapper{
		factory:     factory,
		serviceAddr: serviceAddr,
		client:      factory.GetWebsocketApiGatewayClient(serviceAddr),
		timeout:     500 * time.Millisecond,
	}
}

func (w *WebsocketServiceWrapper) ensureClient() bool {
	if w.client != nil {
		return true
	}

	logger.Log.Info("WebsocketService client is nil, attempting to reconnect...")
	w.client = w.factory.GetWebsocketApiGatewayClient(w.serviceAddr)

	if w.client == nil {
		logger.Log.Warn("Failed to reconnect to WebsocketService")
		return false
	}

	logger.Log.Info("Successfully reconnected to WebsocketService")
	return true
}

func (w *WebsocketServiceWrapper) GetOnline(ctx context.Context) (*websocketpb.GetOnlineResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient() {
		logger.Log.Warn("Using default value for GetOnline due to missing client")
		return &websocketpb.GetOnlineResponse{Online: 0}, nil
	}

	response, err := w.client.GetOnline(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Log.Warnf("GetOnline request failed: %v", err)
		return &websocketpb.GetOnlineResponse{Online: 0}, err
	}

	return response, nil
}

func (w *WebsocketServiceWrapper) GetOnlineSoft(ctx context.Context) *websocketpb.GetOnlineResponse {
	res, err := w.GetOnline(ctx)

	if err != nil {
		return &websocketpb.GetOnlineResponse{Online: 0}
	}

	return res
}

func (w *WebsocketServiceWrapper) GetOnlineUsers(ctx context.Context) (*websocketpb.GetOnlineUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient() {
		logger.Log.Warn("Using default value for GetOnlineUsers due to missing client")
		return &websocketpb.GetOnlineUsersResponse{Users: []*websocketpb.OnlineUser{}}, nil
	}

	response, err := w.client.GetOnlineUsers(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Log.Warnf("GetOnlineUsers request failed: %v", err)
		return &websocketpb.GetOnlineUsersResponse{Users: []*websocketpb.OnlineUser{}}, err
	}

	return response, nil
}

func (w *WebsocketServiceWrapper) GetOnlineUsersSoft(ctx context.Context) *websocketpb.GetOnlineUsersResponse {
	res, err := w.GetOnlineUsers(ctx)

	if err != nil {
		return &websocketpb.GetOnlineUsersResponse{Users: []*websocketpb.OnlineUser{}}
	}

	return res
}

func (w *WebsocketServiceWrapper) SendMessage(ctx context.Context, request *websocketpb.SendMessageRequest) error {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient() {
		logger.Log.Warn("Failed to send message due to missing client")
		return nil
	}

	_, err := w.client.SendMessage(ctx, request)
	if err != nil {
		logger.Log.Warnf("SendMessage request failed: %v", err)
		return err
	}

	return nil
}

func (w *WebsocketServiceWrapper) Broadcast(ctx context.Context, request *websocketpb.BroadcastRequest) error {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient() {
		logger.Log.Warn("Failed to broadcast message due to missing client")
		return nil
	}

	_, err := w.client.Broadcast(ctx, request)
	if err != nil {
		logger.Log.Warnf("Broadcast request failed: %v", err)
		return err
	}

	return nil
}
