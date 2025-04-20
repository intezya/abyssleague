package routes

import (
	"abysscore/internal/adapters/config"
	"abysscore/internal/adapters/controller/http/handlers"
	domainservice "abysscore/internal/domain/service"
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
	"path"
)

type DependencyProvider struct {
	Config                *config.Config
	RedisClient           *rediswrapper.ClientWrapper
	AuthenticationService domainservice.AuthenticationService
	apiPrefix             string

	// Routes map is replaced with RouteGroups
	routeGroups []*RouteGroup

	// Keep the Routes map for backward compatibility if needed
	Routes map[string]*Route
}

func NewDependencyProvider(
	dependencyProvider *handlers.DependencyProvider,
	authenticationService domainservice.AuthenticationService,
	redisClient *rediswrapper.ClientWrapper,
	config *config.Config,
) *DependencyProvider {
	apiPrefix := "/api"
	dp := &DependencyProvider{
		Config:                config,
		RedisClient:           redisClient,
		AuthenticationService: authenticationService,
		apiPrefix:             apiPrefix,
		Routes:                make(map[string]*Route),
	}

	// Setup route groups
	dp.setupRouteGroups(dependencyProvider)

	// For backward compatibility, convert route groups to the flat map
	dp.populateRoutesMap()

	return dp
}

// setupRouteGroups organizes routes into logical groups
func (dp *DependencyProvider) setupRouteGroups(handlers *handlers.DependencyProvider) {
	authGroup := GetAuthGroup(handlers, dp)
	gameItemGroup := GetGameItemGroup(handlers, dp)

	dp.routeGroups = []*RouteGroup{authGroup, gameItemGroup}
}

// WithoutAuthenticationRequirement disables authentication for a route
func WithoutAuthenticationRequirement() RouteOption {
	return func(r *Route) {
		r.RequireAuthentication = false
	}
}

// populateRoutesMap converts route groups to the flat map for backward compatibility
func (dp *DependencyProvider) populateRoutesMap() {
	for _, group := range dp.routeGroups {
		for _, entry := range group.Routes {
			fullPath := path.Join(group.Prefix, entry.Path)
			dp.Routes[fullPath] = entry.Route
		}
	}
}

// GetRouteGroups returns the route groups for direct use
func (dp *DependencyProvider) GetRouteGroups() []*RouteGroup {
	return dp.routeGroups
}
