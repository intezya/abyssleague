package wrapper

import (
	"context"
	"errors"
	"sync"
	"time"

	websocketpb "github.com/intezya/abyssleague/proto/websocket"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/factory"
	"github.com/intezya/pkglib/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

var errInvalidClientType = errors.New("invalid client type")

const (
	defaultGRPCTimeout = 1 * time.Second
	reconnectCooldown  = 5 * time.Second
)

type WebsocketServiceWrapper struct {
	factory     *factory.GrpcClientFactory
	serviceAddr string
	client      websocketpb.WebsocketServiceClient
	timeout     time.Duration

	// Added fields for connection management
	reconnectMu       sync.Mutex
	lastReconnectTime time.Time
}

func NewWebsocketServiceWrapper(
	ctx context.Context,
	factory *factory.GrpcClientFactory,
	serviceAddr string,
) *WebsocketServiceWrapper {
	wrapper := &WebsocketServiceWrapper{
		factory:     factory,
		serviceAddr: serviceAddr,
		timeout:     defaultGRPCTimeout,
		client:      nil,
	}

	// Initialize client immediately, but don't block
	go func() {
		_, err := factory.GetAndSetWebsocketApiGatewayClient(ctx, serviceAddr, wrapper)
		if err != nil {
			logger.Log.Warnf("Failed to set WebsocketApiGateway client: %v", err)
		}
	}()

	return wrapper
}

func (w *WebsocketServiceWrapper) SetClient(client interface{}) error {
	if typedClient, ok := client.(websocketpb.WebsocketServiceClient); ok {
		w.client = typedClient

		return nil
	}

	return errInvalidClientType
}

func (w *WebsocketServiceWrapper) GetOnline(
	ctx context.Context,
) (*websocketpb.GetOnlineResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient(ctx) {
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

func (w *WebsocketServiceWrapper) GetOnlineUsers(
	ctx context.Context,
) (*websocketpb.GetOnlineUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient(ctx) {
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

func (w *WebsocketServiceWrapper) SendMessage(
	ctx context.Context,
	request *websocketpb.SendMessageRequest,
) error {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient(ctx) {
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

func (w *WebsocketServiceWrapper) Broadcast(
	ctx context.Context,
	request *websocketpb.BroadcastRequest,
) error {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	if !w.ensureClient(ctx) {
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

// Optimized client reconnection logic with rate limiting.
func (w *WebsocketServiceWrapper) ensureClient(ctx context.Context) bool {
	// If we already have a client, return immediately
	if w.client != nil {
		return true
	}

	// Lock to prevent multiple goroutines from attempting reconnection simultaneously
	w.reconnectMu.Lock()
	defer w.reconnectMu.Unlock()

	// Check again after acquiring lock (double-checked locking pattern)
	if w.client != nil {
		return true
	}

	// Implement reconnection rate limiting
	now := time.Now()
	if now.Sub(w.lastReconnectTime) < reconnectCooldown {
		logger.Log.Debug("Skipping reconnection attempt due to cooldown period")

		return false
	}

	// Update last reconnect timestamp
	w.lastReconnectTime = now

	logger.Log.Info("WebsocketService client is nil, attempting to reconnect...")

	client, err := w.factory.GetAndSetWebsocketApiGatewayClient(ctx, w.serviceAddr, nil)
	if err != nil {
		logger.Log.Warnf("Failed to reconnect to WebsocketService: %v", err)

		return false
	}

	w.client = client

	logger.Log.Info("Successfully reconnected to WebsocketService")

	return true
}
