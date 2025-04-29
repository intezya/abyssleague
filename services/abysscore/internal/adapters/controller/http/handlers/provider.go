package handlers

import (
	applicationservice "github.com/intezya/abyssleague/services/abysscore/internal/application/service"
)

type DependencyProvider struct {
	appProvider *applicationservice.DependencyProvider

	AuthenticationHandler *AuthenticationHandler
	GameItemHandler       *GameItemHandler
	InventoryItemHandler  *InventoryItemHandler
	AccountHandler        *AccountHandler
}

func NewDependencyProvider(
	dependencyProvider *applicationservice.DependencyProvider,
) *DependencyProvider {
	return &DependencyProvider{
		appProvider: dependencyProvider,

		AuthenticationHandler: NewAuthenticationHandler(dependencyProvider.AuthenticationService),
		GameItemHandler:       NewGameItemHandler(dependencyProvider.GameItemService),
		InventoryItemHandler:  NewInventoryItemHandler(dependencyProvider.InventoryItemService),
		AccountHandler:        NewAccountHandler(dependencyProvider.AccountService),
	}
}
