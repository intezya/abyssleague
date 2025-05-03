package clients

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/grpcwrap"
)

type DependencyProvider struct {
	MainWebsocketService  WebsocketMessagingClient
	DraftWebsocketService WebsocketMessagingClient
}

func NewDependencyProvider(
	gRPCConfig *Config,
) *DependencyProvider {
	return &DependencyProvider{
		NewWebsocketMessagingClient(
			gRPCConfig.MainWebsocketServerAddress(),
			grpcwrap.WithDevMode(gRPCConfig.DevMode),
			grpcwrap.WithRPCTimeout(gRPCConfig.RequestTimeout),
		),
		NewWebsocketMessagingClient(
			gRPCConfig.DraftWebsocketServerAddress(),
			grpcwrap.WithDevMode(gRPCConfig.DevMode),
			grpcwrap.WithRPCTimeout(gRPCConfig.RequestTimeout),
		),
	}
}

func (d *DependencyProvider) CloseAll() error {
	err := d.MainWebsocketService.Close()
	if err != nil {
		return err
	}

	err = d.DraftWebsocketService.Close()
	if err != nil {
		return err
	}

	return nil
}
