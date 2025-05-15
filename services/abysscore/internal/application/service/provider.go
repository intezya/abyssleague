package applicationservice

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/clients"
	drivenports "github.com/intezya/abyssleague/services/abysscore/internal/domain/ports/driven"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	eventlib "github.com/intezya/abyssleague/services/abysscore/pkg/event"
)

type DependencyProvider struct {
	repositoryDependencyProvider *persistence.DependencyProvider
	gRPCDependencyProvider       *clients.DependencyProvider
	passwordHelper               domainservice.CredentialsHelper
	tokenHelper                  domainservice.TokenHelper

	EventPublisher eventlib.ApplicationEventPublisher

	AuthenticationService domainservice.AuthenticationService
	GameItemService       domainservice.GameItemService
	InventoryItemService  domainservice.InventoryItemService
	AccountService        domainservice.AccountService
}

func NewDependencyProvider(
	repositoryDependencyProvider *persistence.DependencyProvider,
	gRPCDependencyProvider *clients.DependencyProvider,
	passwordHelper domainservice.CredentialsHelper,
	tokenHelper domainservice.TokenHelper,
	mailSender drivenports.MailSender,
) *DependencyProvider {
	mainClientNotificationService := NewNotificationService(
		gRPCDependencyProvider.MainWebsocketService,
	)
	// draftClientNotificationService := NewNotificationService(gRPCDependencyProvider.DraftWebsocketService)
	return &DependencyProvider{
		repositoryDependencyProvider: repositoryDependencyProvider,
		gRPCDependencyProvider:       gRPCDependencyProvider,
		passwordHelper:               passwordHelper,
		tokenHelper:                  tokenHelper,

		AuthenticationService: NewAuthenticationService(
			repositoryDependencyProvider.AuthenticationRepository,
			repositoryDependencyProvider.UserRepository,
			gRPCDependencyProvider.MainWebsocketService,
			passwordHelper,
			tokenHelper,
			repositoryDependencyProvider.BannedHardwareIDRepository,
			NewAuthenticationEventService(
				repositoryDependencyProvider.UserRepository,
			),
		),
		GameItemService: NewGameItemService(repositoryDependencyProvider.GameItemRepository),
		InventoryItemService: NewInventoryItemService(
			repositoryDependencyProvider.InventoryItemRepository,
			repositoryDependencyProvider.InventoryRepository,
			NewInventoryItemEventService(
				mainClientNotificationService,
			),
		),
		AccountService: NewAccountService(
			repositoryDependencyProvider.UserRepository,
			mailSender,
			repositoryDependencyProvider.MailMessageRepository,
		),
	}
}
