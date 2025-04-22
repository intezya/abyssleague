package factory

import (
	websocketpb "abyssproto/websocket"
	"fmt"
	"github.com/intezya/pkglib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

type GRPCConfig struct {
	WebsocketApiGatewayHost  string
	WebsocketApiGatewayPorts []int
}

const mainWebsocketServerIdx = 0
const draftWebsocketServerIdx = 1

func (g *GRPCConfig) MainWebsocketServerAddress() string {
	if len(g.WebsocketApiGatewayHost) == mainWebsocketServerIdx {
		return ""
	}
	return fmt.Sprintf("%s:%d", g.WebsocketApiGatewayHost, g.WebsocketApiGatewayPorts[mainWebsocketServerIdx])
}

func (g *GRPCConfig) DraftWebsocketServerAddress() string {
	if len(g.WebsocketApiGatewayHost) <= draftWebsocketServerIdx {
		return ""
	}
	return fmt.Sprintf("%s:%d", g.WebsocketApiGatewayHost, g.WebsocketApiGatewayPorts[draftWebsocketServerIdx])
}

type GrpcClientFactory struct {
	connections map[string]*grpc.ClientConn
	clients     map[string]interface{}
}

func NewGrpcClientFactory() *GrpcClientFactory {
	return &GrpcClientFactory{
		connections: make(map[string]*grpc.ClientConn),
		clients:     make(map[string]interface{}),
	}
}

func (f *GrpcClientFactory) GetWebsocketApiGatewayClient(address string) websocketpb.WebsocketServiceClient {
	key := "websocket-" + address

	if client, exists := f.clients[key]; exists {
		logger.Log.Info("Using existing connection: ", address)
		return client.(websocketpb.WebsocketServiceClient)
	}

	maxRetries := 1
	retryInterval := 1 * time.Second
	var conn *grpc.ClientConn
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logger.Log.Infof("Connection attempt to %s (%d of %d)", address, attempt, maxRetries)

		conn, err = grpc.Dial(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(10*1024*1024),
				grpc.MaxCallSendMsgSize(10*1024*1024),
			),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                20 * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			}),
		)

		if err == nil {
			logger.Log.Infof("Successfully connected to %s (%d attempmt)", address, attempt)
			break
		}

		logger.Log.Infof("Connection to %s failed: %v , retrying in %v seconds...", address, err, retryInterval)

		if attempt < maxRetries {
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		logger.Log.Warn("Failed to connect to Websocket service: ", err)
		return nil
	}

	f.connections[key] = conn
	client := websocketpb.NewWebsocketServiceClient(conn)
	f.clients[key] = client

	return client
}

func (f *GrpcClientFactory) CloseAll() {
	for _, conn := range f.connections {
		_ = conn.Close()
	}
}
