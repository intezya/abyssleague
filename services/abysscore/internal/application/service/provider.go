package applicationservice

import (
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/infrastructure/persistence"
)

type DependencyProvider struct {
	dependencyProvider *persistence.DependencyProvider
	passwordHelper     domainservice.CredentialsHelper
	tokenHelper        domainservice.TokenHelper

	AuthenticationService domainservice.AuthenticationService
}

func NewDependencyProvider(
	dependencyProvider *persistence.DependencyProvider,
	passwordHelper domainservice.CredentialsHelper,
	tokenHelper domainservice.TokenHelper,
) *DependencyProvider {
	return &DependencyProvider{
		dependencyProvider: dependencyProvider,
		passwordHelper:     passwordHelper,
		tokenHelper:        tokenHelper,

		AuthenticationService: NewAuthenticationService(dependencyProvider.UserRepository, passwordHelper, tokenHelper),
	}
}
