package factory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	websocketpb "github.com/intezya/abyssleague/proto/websocket"
	"github.com/intezya/pkglib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Server indices for websocket servers.
const (
	MainWebsocketServerIdx  = 0
	DraftWebsocketServerIdx = 1
)

// GRPC connection parameters.
const (
	GRPCMaxRetries     = 10
	GRPCReconnectDelay = 5 * time.Second

	GRPCMaxCallRecvMsgSize = 10 * 1024 * 1024
	GRPCMaxCallSendMsgSize = 10 * 1024 * 1024

	GRPCKeepAliveTime    = 20 * time.Second
	GRPCKeepAliveTimeout = 5 * time.Second
)

var (
	errEmptyAddress        = errors.New("cannot connect to empty address")
	errNoClientExists      = errors.New("no client exists")
	errConnectionMissing   = errors.New("client exists but connection missing")
	errInvalidClientType   = errors.New("invalid client type")
	errUnhealthyConnection = errors.New("unhealthy connection")
)

// GRPCConfig holds configuration for GRPC connections.
type GRPCConfig struct {
	WebsocketApiGatewayHost  string
	WebsocketApiGatewayPorts []int
}

// MainWebsocketServerAddress returns the address of the main websocket server.
func (g *GRPCConfig) MainWebsocketServerAddress() string {
	if g.WebsocketApiGatewayHost == "" ||
		len(g.WebsocketApiGatewayPorts) <= MainWebsocketServerIdx {
		return ""
	}

	return fmt.Sprintf(
		"%s:%d",
		g.WebsocketApiGatewayHost,
		g.WebsocketApiGatewayPorts[MainWebsocketServerIdx],
	)
}

// DraftWebsocketServerAddress returns the address of the draft websocket server.
func (g *GRPCConfig) DraftWebsocketServerAddress() string {
	if g.WebsocketApiGatewayHost == "" ||
		len(g.WebsocketApiGatewayPorts) <= DraftWebsocketServerIdx {
		return ""
	}

	return fmt.Sprintf(
		"%s:%d",
		g.WebsocketApiGatewayHost,
		g.WebsocketApiGatewayPorts[DraftWebsocketServerIdx],
	)
}

// ClientReceiver interface for components that need to receive client references.
type ClientReceiver interface {
	SetClient(client interface{}) error
}

// GrpcClientFactory manages GRPC client connections.
type GrpcClientFactory struct {
	connections map[string]*grpc.ClientConn
	clients     map[string]interface{}
	mu          sync.RWMutex
}

// NewGrpcClientFactory creates a new instance of GrpcClientFactory.
func NewGrpcClientFactory() *GrpcClientFactory {
	return &GrpcClientFactory{
		connections: make(map[string]*grpc.ClientConn),
		clients:     make(map[string]interface{}),
	}
}

// GetAndSetWebsocketApiGatewayClient returns a websocket service client for the given address.
func (f *GrpcClientFactory) GetAndSetWebsocketApiGatewayClient(
	ctx context.Context,
	address string,
	receiver ClientReceiver,
) (websocketpb.WebsocketServiceClient, error) {
	if address == "" {
		return nil, errEmptyAddress
	}

	key := "websocket-" + address

	// Try to get existing client
	client, err := f.getExistingClient(key)
	if err == nil {
		// Set client on receiver if provided
		if receiver != nil {
			if err := receiver.SetClient(client); err != nil {
				return nil, fmt.Errorf("failed to set client on receiver: %w", err)
			}
		}

		return client, nil
	}

	// Create new client
	client, err = f.createNewClient(ctx, address, key)
	if err != nil {
		return nil, err
	}

	// Set client on receiver if provided
	if receiver != nil {
		if err := receiver.SetClient(client); err != nil {
			return nil, fmt.Errorf("failed to set client on receiver: %w", err)
		}
	}

	return client, nil
}

// CloseAll closes all connections managed by the factory.
func (f *GrpcClientFactory) CloseAll() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for addr, conn := range f.connections {
		logger.Log.Infof("Closing GRPC connection: %s", addr)

		if err := conn.Close(); err != nil {
			logger.Log.Warnf("Error closing connection to %s: %v", addr, err)
		}
	}

	// Clear maps
	f.connections = make(map[string]*grpc.ClientConn)
	f.clients = make(map[string]interface{})
}

// GetConnectionCount returns the number of active connections.
func (f *GrpcClientFactory) GetConnectionCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return len(f.connections)
}

// IsConnected checks if there's an active connection to the given address.
func (f *GrpcClientFactory) IsConnected(address string) bool {
	if address == "" {
		return false
	}

	key := "websocket-" + address

	f.mu.RLock()
	defer f.mu.RUnlock()

	conn, exists := f.connections[key]
	if !exists {
		return false
	}

	state := conn.GetState()

	return state == connectivity.Ready || state == connectivity.Idle
}

// getExistingClient tries to return an existing healthy client.
func (f *GrpcClientFactory) getExistingClient(
	key string,
) (websocketpb.WebsocketServiceClient, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	client, exists := f.clients[key]
	if !exists {
		return nil, errNoClientExists
	}

	logger.Log.Info("Using existing GRPC connection:", key)

	// Check connection health
	conn, ok := f.connections[key]
	if !ok {
		return nil, errConnectionMissing
	}

	state := conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle {
		return nil, fmt.Errorf("%w: %s", errUnhealthyConnection, state)
	}

	// Return typed client
	typedClient, ok := client.(websocketpb.WebsocketServiceClient)
	if !ok {
		return nil, errInvalidClientType
	}

	return typedClient, nil
}

// createNewClient creates and stores a new client.
func (f *GrpcClientFactory) createNewClient(
	ctx context.Context,
	address string,
	key string,
) (websocketpb.WebsocketServiceClient, error) {
	// Create connection
	conn, err := f.connectWithRetry(ctx, address)
	if err != nil {
		return nil, err
	}

	// Create client from connection
	client := websocketpb.NewWebsocketServiceClient(conn)

	// Store connection and client
	f.mu.Lock()
	defer f.mu.Unlock()

	f.connections[key] = conn
	f.clients[key] = client

	return client, nil
}

// connectWithRetry attempts to establish a GRPC connection with retry logic.
func (f *GrpcClientFactory) connectWithRetry(
	ctx context.Context,
	address string,
) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn

	var err error

	for attempt := 1; attempt <= GRPCMaxRetries; attempt++ {
		logger.Log.Infof(
			"GRPC connection attempt to %s (%d of %d)",
			address,
			attempt,
			GRPCMaxRetries,
		)

		// Create connection options
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(GRPCMaxCallRecvMsgSize),
				grpc.MaxCallSendMsgSize(GRPCMaxCallSendMsgSize),
			),
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:                GRPCKeepAliveTime,
					Timeout:             GRPCKeepAliveTimeout,
					PermitWithoutStream: true,
				},
			),
		}

		conn, err = grpc.NewClient(address, opts...)

		if err == nil {
			logger.Log.Infof("Successfully connected to %s (attempt %d)", address, attempt)

			return conn, nil
		}

		logger.Log.Warnf(
			"Connection to %s failed: %v, retrying in %v...",
			address,
			err,
			GRPCReconnectDelay,
		)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf(
				"context canceled while connecting to %s: %w",
				address,
				ctx.Err(),
			)
		case <-time.After(GRPCReconnectDelay): // Continue to next attempt
		}
	}

	return nil, fmt.Errorf(
		"failed to connect to %s after %d attempts: %w",
		address,
		GRPCMaxRetries,
		err,
	)
}
