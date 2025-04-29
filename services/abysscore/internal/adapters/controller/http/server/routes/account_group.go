package routes

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/handlers"
	"path"
)

func GetAccountGroup(
	handlers *handlers.DependencyProvider,
	provider *DependencyProvider,
) *RouteGroup {
	accountGroup := NewRouteGroup(path.Join(provider.apiPrefix, "users"))

	accountGroup.Add(
		"/account/email/get_code", NewRoute(
			handlers.AccountHandler.SendCodeForEmailLink,
			MethodPost,
		),
	)

	accountGroup.Add(
		"/account/email/enter_code",
		NewRoute(
			handlers.AccountHandler.EnterCodeForEmailLink,
			MethodPost,
		),
	)

	return accountGroup
}
