package server

import (
	"abysscore/internal/adapters/config"
	"abysscore/internal/adapters/controller/http/handlers"
	domainservice "abysscore/internal/domain/service"
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
	"path"
)

type DependencyProvider struct {
	config                *config.Config
	redisClient           *rediswrapper.ClientWrapper
	authenticationService domainservice.AuthenticationService
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
		config:                config,
		redisClient:           redisClient,
		authenticationService: authenticationService,
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
	// Authentication routes group
	authGroup := NewRouteGroup(path.Join(dp.apiPrefix, "auth"))

	// Add routes to the auth group
	authGroup.Add(
		"/register", NewRoute(
			handlers.AuthenticationHandler.Register,
			MethodPost,
			WithoutAuthenticationRequirement(),
			WithRateLimit(AuthRateLimit),
		),
	)

	authGroup.Add(
		"/login", NewRoute(
			handlers.AuthenticationHandler.Login, // Fixed: was using Register here
			MethodPost,
			WithoutAuthenticationRequirement(),
			WithRateLimit(AuthRateLimit),
		),
	)

	dp.routeGroups = []*RouteGroup{authGroup}
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
		for _, entry := range group.routes {
			fullPath := path.Join(group.prefix, entry.path)
			dp.Routes[fullPath] = entry.route
		}
	}
}

// GetRouteGroups returns the route groups for direct use
func (dp *DependencyProvider) GetRouteGroups() []*RouteGroup {
	return dp.routeGroups
}
