package wrapper

import (
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/factory"
)

const (
	defaultGRPCTimeout = 500 * time.Millisecond
)

type DependencyProvider struct {
	config *factory.GRPCConfig

	MainWebsocketService  *WebsocketServiceWrapper
	DraftWebsocketService *WebsocketServiceWrapper
}

func NewDependencyProvider(
	config *factory.GRPCConfig,
	factory *factory.GrpcClientFactory,
) *DependencyProvider {
	return &DependencyProvider{
		config: config,
		MainWebsocketService: NewWebsocketServiceWrapper(
			factory,
			config.MainWebsocketServerAddress(),
		),
		DraftWebsocketService: NewWebsocketServiceWrapper(
			factory,
			config.DraftWebsocketServerAddress(),
		),
	}
}
