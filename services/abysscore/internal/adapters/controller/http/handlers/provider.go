package handlers

import (
	applicationservice "abysscore/internal/application/service"
	"abysscore/internal/pkg/validator"
)

type DependencyProvider struct {
	dependencyProvider *applicationservice.DependencyProvider
	v                  *validator.Validator

	AuthenticationHandler *Authentication
}

func NewDependencyProvider(
	dependencyProvider *applicationservice.DependencyProvider,
	v *validator.Validator,
) *DependencyProvider {
	return &DependencyProvider{
		dependencyProvider:    dependencyProvider,
		v:                     v,
		AuthenticationHandler: NewAuthenticationHandler(dependencyProvider.AuthenticationService, v),
	}
}
