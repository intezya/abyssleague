package applicationservice

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/wrapper"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	eventlib "github.com/intezya/abyssleague/services/abysscore/pkg/event"
)

type DependencyProvider struct {
	repositoryDependencyProvider *persistence.DependencyProvider
	gRPCDependencyProvider       *wrapper.DependencyProvider
	passwordHelper               domainservice.CredentialsHelper
	tokenHelper                  domainservice.TokenHelper

	EventPublisher eventlib.Publisher

	AuthenticationService domainservice.AuthenticationService
	GameItemService       domainservice.GameItemService
	InventoryItemService  domainservice.InventoryItemService
}

func NewDependencyProvider(
	repositoryDependencyProvider *persistence.DependencyProvider,
	gRPCDependencyProvider *wrapper.DependencyProvider,
	passwordHelper domainservice.CredentialsHelper,
	tokenHelper domainservice.TokenHelper,
) *DependencyProvider {
	mainClientNotificationService := NewNotificationService(gRPCDependencyProvider.MainWebsocketService)
	//draftClientNotificationService := NewNotificationService(gRPCDependencyProvider.DraftWebsocketService)

	eventPublisher := NewApplicationEventPublisher(mainClientNotificationService)

	return &DependencyProvider{
		repositoryDependencyProvider: repositoryDependencyProvider,
		gRPCDependencyProvider:       gRPCDependencyProvider,
		passwordHelper:               passwordHelper,
		tokenHelper:                  tokenHelper,

		EventPublisher: eventPublisher,

		AuthenticationService: NewAuthenticationService(
			repositoryDependencyProvider.UserRepository,
			gRPCDependencyProvider.MainWebsocketService,
			passwordHelper,
			tokenHelper,
		),
		GameItemService: NewGameItemService(repositoryDependencyProvider.GameItemRepository),
		InventoryItemService: NewInventoryItemService(
			repositoryDependencyProvider.InventoryItemRepository,
			repositoryDependencyProvider.UserRepository,
			eventPublisher,
		),
	}
}
