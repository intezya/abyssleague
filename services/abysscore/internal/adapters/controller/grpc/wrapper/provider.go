package wrapper

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/factory"
)

type DependencyProvider struct {
	config *factory.GRPCConfig

	MainWebsocketService  *WebsocketServiceWrapper
	DraftWebsocketService *WebsocketServiceWrapper
}

func NewDependencyProvider(
	ctx context.Context,
	config *factory.GRPCConfig,
	factory *factory.GrpcClientFactory,
) *DependencyProvider {
	return &DependencyProvider{
		config: config,
		MainWebsocketService: NewWebsocketServiceWrapper(
			ctx,
			factory,
			config.MainWebsocketServerAddress(),
		),
		DraftWebsocketService: NewWebsocketServiceWrapper(
			ctx,
			factory,
			config.DraftWebsocketServerAddress(),
		),
	}
}
