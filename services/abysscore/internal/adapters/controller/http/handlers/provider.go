package handlers

import (
	applicationservice "abysscore/internal/application/service"
)

type DependencyProvider struct {
	appProvider *applicationservice.DependencyProvider

	AuthenticationHandler *AuthenticationHandler
	GameItemHandler       *GameItemHandler
}

func NewDependencyProvider(
	dependencyProvider *applicationservice.DependencyProvider,
) *DependencyProvider {
	return &DependencyProvider{
		appProvider: dependencyProvider,

		AuthenticationHandler: NewAuthenticationHandler(dependencyProvider.AuthenticationService),
		GameItemHandler:       NewGameItemHandler(dependencyProvider.GameItemService),
	}
}
