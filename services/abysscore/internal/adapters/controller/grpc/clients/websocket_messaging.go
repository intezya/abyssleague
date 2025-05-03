package clients

import (
	"context"
	"time"

	websocketpb "github.com/intezya/abyssleague/proto/websocket"
	"github.com/intezya/abyssleague/services/abysscore/pkg/grpcwrap"
	"github.com/intezya/pkglib/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// OnlineUser - custom user model.
type OnlineUser struct {
	ID         int
	Username   string
	HardwareID string
}

// WebsocketMessagingClient custom client interface.
type WebsocketMessagingClient interface {
	GetOnline(ctx context.Context) (int, error)
	GetOnlineUsers(ctx context.Context) ([]*OnlineUser, error)
	SendMessage(ctx context.Context, userID int, jsonPayload []byte) error
	Broadcast(ctx context.Context, jsonPayload []byte) error
	WaitForConnection(ctx context.Context) error
	Close() error
}

type websocketMessagingClientImpl struct {
	grpcClient *grpcwrap.GenericGRPCClient[websocketpb.WebsocketServiceClient]

	// Converters
	usersConverter grpcwrap.TypeConverter[[]*websocketpb.OnlineUser, []*OnlineUser]
}

// Checking, that we implement WebsocketMessagingClient interface.
var _ WebsocketMessagingClient = (*websocketMessagingClientImpl)(nil)

func NewWebsocketMessagingClient(
	serviceAddr string,
	opts ...grpcwrap.ClientOption,
) WebsocketMessagingClient {
	clientCreator := func(conn *grpc.ClientConn) websocketpb.WebsocketServiceClient {
		return websocketpb.NewWebsocketServiceClient(conn)
	}

	baseClient := grpcwrap.NewBaseGRPCClient(serviceAddr, clientCreator, opts...)

	userConverter := &grpcwrap.SimpleConverter[*websocketpb.OnlineUser, *OnlineUser]{
		ConvertFunc: func(from *websocketpb.OnlineUser) (*OnlineUser, error) {
			return &OnlineUser{
				ID:         int(from.GetId()),
				Username:   from.GetUsername(),
				HardwareID: from.GetHardwareID(),
			}, nil
		},
	}

	usersConverter := &grpcwrap.SliceConverter[*websocketpb.OnlineUser, *OnlineUser]{
		ElementConverter: userConverter,
	}

	client := &websocketMessagingClientImpl{
		grpcClient:     baseClient,
		usersConverter: usersConverter,
	}

	// Warm connection
	go func(client *websocketMessagingClientImpl) {
		for {
			time.Sleep(time.Second * 1)

			if client.grpcClient.BaseGRPCClient == nil {
				continue
			}

			if client.grpcClient.ConnectionWarm {
				break
			}

			_, _ = client.GetOnline(context.Background())
			client.grpcClient.ConnectionWarm = true

			logger.Log.Debugf("WebsocketMessagingClient (%s) Client warmed!", serviceAddr)
		}
	}(client)

	return client
}

func (c *websocketMessagingClientImpl) GetOnline(ctx context.Context) (int, error) {
	resp, err := grpcwrap.ExecuteCallWithFallback[
		websocketpb.WebsocketServiceClient,
		*emptypb.Empty,
		*websocketpb.GetOnlineResponse,
	](
		c.grpcClient,
		ctx,
		func(client websocketpb.WebsocketServiceClient, req *emptypb.Empty) (*websocketpb.GetOnlineResponse, error) {
			return client.GetOnline(ctx, req)
		},
		&emptypb.Empty{},
		&websocketpb.GetOnlineResponse{Online: 0},
	)
	if err != nil {
		return 0, err
	}

	return int(resp.GetOnline()), nil
}

func (c *websocketMessagingClientImpl) GetOnlineUsers(ctx context.Context) ([]*OnlineUser, error) {
	resp, err := grpcwrap.ExecuteCallWithFallback[
		websocketpb.WebsocketServiceClient,
		*emptypb.Empty,
		*websocketpb.GetOnlineUsersResponse,
	](
		c.grpcClient,
		ctx,
		func(
			client websocketpb.WebsocketServiceClient,
			req *emptypb.Empty,
		) (*websocketpb.GetOnlineUsersResponse, error) {
			return client.GetOnlineUsers(ctx, req)
		},
		&emptypb.Empty{},
		&websocketpb.GetOnlineUsersResponse{Users: []*websocketpb.OnlineUser{}},
	)
	if err != nil {
		return nil, err
	}

	users, err := c.usersConverter.Convert(resp.GetUsers())
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *websocketMessagingClientImpl) SendMessage(
	ctx context.Context,
	userID int,
	jsonPayload []byte,
) error {
	client, err := c.grpcClient.GetClient()
	if err != nil {
		if c.grpcClient.DevMode {
			return nil
		}

		return err
	}

	_, err = client.SendMessage(ctx, &websocketpb.SendMessageRequest{
		UserId:      int64(userID),
		JsonPayload: jsonPayload,
	})
	if err != nil && c.grpcClient.DevMode {
		return nil
	}

	return err
}

func (c *websocketMessagingClientImpl) Broadcast(ctx context.Context, jsonPayload []byte) error {
	client, err := c.grpcClient.GetClient()
	if err != nil {
		if c.grpcClient.DevMode {
			return nil
		}

		return err
	}

	_, err = client.Broadcast(ctx, &websocketpb.BroadcastRequest{
		JsonPayload: jsonPayload,
	})
	if err != nil && c.grpcClient.DevMode {
		return nil
	}

	return err
}

func (c *websocketMessagingClientImpl) WaitForConnection(ctx context.Context) error {
	return c.grpcClient.WaitForConnection(ctx)
}

func (c *websocketMessagingClientImpl) Close() error {
	return c.grpcClient.Close()
}
