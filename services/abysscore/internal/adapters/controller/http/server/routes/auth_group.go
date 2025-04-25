package routes

import (
	"path"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/handlers"
)

func GetAuthGroup(handlers *handlers.DependencyProvider, provider *DependencyProvider) *RouteGroup {
	authGroup := NewRouteGroup(path.Join(provider.apiPrefix, "auth"))

	authGroup.Add(
		"/register", NewRoute(
			handlers.AuthenticationHandler.Register,
			MethodPost,
			WithoutAuthenticationRequirement(),
			WithRateLimit(AuthRateLimit),
		),
	)

	authGroup.Add(
		"/login",
		NewRoute(
			handlers.AuthenticationHandler.Login,
			MethodPost,
			WithoutAuthenticationRequirement(),
			WithRateLimit(AuthRateLimit),
		),
	)

	authGroup.Add(
		"/change_password",
		NewRoute(
			handlers.AuthenticationHandler.ChangePassword,
			MethodPost,
			WithoutAuthenticationRequirement(),
			WithRateLimit(AuthRateLimit),
		),
	)

	return authGroup
}
