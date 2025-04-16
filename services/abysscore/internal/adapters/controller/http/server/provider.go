package server

import (
	"abysscore/internal/adapters/config"
	"abysscore/internal/adapters/controller/http/handlers"
	domainservice "abysscore/internal/domain/service"
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
)

type DependencyProvider struct {
	redisClient           *rediswrapper.ClientWrapper
	config                *config.Config
	handlerDependencies   *handlers.DependencyProvider
	authenticationService domainservice.AuthenticationService
	authenticationHandler *handlers.Authentication
}

func NewDependencyProvider(
	redis *rediswrapper.ClientWrapper,
	config *config.Config,
	handlerDependencies *handlers.DependencyProvider,
	authenticationService domainservice.AuthenticationService,
) *DependencyProvider {
	return &DependencyProvider{
		redisClient:           redis,
		config:                config,
		handlerDependencies:   handlerDependencies,
		authenticationService: authenticationService,
		authenticationHandler: handlerDependencies.AuthenticationHandler,
	}
}
