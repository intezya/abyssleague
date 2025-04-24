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

const (
	gRPCMaxRetries     = 10
	gRPCReconnectDelay = 5 * time.Second

	gRPCMaxCallRecvMsgSize = 10 * 1024 * 1024
	gRPCMaxCallSendMsgSize = 10 * 1024 * 1024

	gRPCKeepAliveTime    = 20 * time.Second // ???
	gRPCKeepAliveTimeout = 5 * time.Second
)

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

type ClientReceiver interface {
	SetClient(client interface{})
}

func (f *GrpcClientFactory) GetAndSetWebsocketApiGatewayClient(
	address string,
	receiver ClientReceiver,
) websocketpb.WebsocketServiceClient {
	key := "websocket-" + address

	if client, exists := f.clients[key]; exists {
		logger.Log.Info("Using existing connection: ", address)

		if receiver != nil {
			receiver.SetClient(client)
		}

		return client.(websocketpb.WebsocketServiceClient)
	}

	var conn *grpc.ClientConn

	var err error

	for attempt := 1; attempt <= gRPCMaxRetries; attempt++ {
		logger.Log.Infof("Connection attempt to %s (%d of %d)", address, attempt, gRPCMaxRetries)

		conn, err = grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(gRPCMaxCallRecvMsgSize),
				grpc.MaxCallSendMsgSize(gRPCMaxCallSendMsgSize),
			),
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:                gRPCKeepAliveTime,
					Timeout:             gRPCKeepAliveTimeout,
					PermitWithoutStream: true,
				},
			),
		)

		if err == nil {
			logger.Log.Infof("Successfully connected to %s (%d attempmt)", address, attempt)

			break
		}

		logger.Log.Infof("Connection to %s failed: %v , retrying in %v seconds...", address, err, gRPCReconnectDelay)

		if attempt < gRPCMaxRetries {
			time.Sleep(gRPCReconnectDelay)
		}
	}

	if err != nil {
		logger.Log.Warn("Failed to connect to Websocket service: ", err)

		return nil
	}

	f.connections[key] = conn
	client := websocketpb.NewWebsocketServiceClient(conn)
	f.clients[key] = client

	if receiver != nil {
		receiver.SetClient(client)
	}

	return client
}

func (f *GrpcClientFactory) CloseAll() {
	for _, conn := range f.connections {
		_ = conn.Close()
	}
}
