package wrapper

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/factory"
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

func (w *WebsocketServiceWrapper) SetClient(client interface{}) {
	w.client = client.(websocketpb.WebsocketServiceClient)
}

func NewWebsocketServiceWrapper(factory *factory.GrpcClientFactory, serviceAddr string) *WebsocketServiceWrapper {
	wrapper := &WebsocketServiceWrapper{
		factory:     factory,
		serviceAddr: serviceAddr,
		timeout:     defaultGRPCTimeout,
	}

	go factory.GetAndSetWebsocketApiGatewayClient(serviceAddr, wrapper)

	return wrapper
}

func (w *WebsocketServiceWrapper) ensureClient() bool {
	if w.client != nil {
		return true
	}

	logger.Log.Info("WebsocketService client is nil, attempting to reconnect...")

	w.client = w.factory.GetAndSetWebsocketApiGatewayClient(w.serviceAddr, nil)

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
