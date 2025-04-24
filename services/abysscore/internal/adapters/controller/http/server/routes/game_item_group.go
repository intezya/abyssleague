package routes

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/handlers"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/schema/access_level"
	"path"
)

func GetGameItemGroup(handlers *handlers.DependencyProvider, provider *DependencyProvider) *RouteGroup {
	gameItemGroup := NewRouteGroup(path.Join(provider.apiPrefix, "items"))

	gameItemGroup.Add(
		"", NewRoute(
			handlers.GameItemHandler.Create,
			MethodPost,
			WithAccessLevel(access_level.CreateItem),
		),
	)

	gameItemGroup.Add(
		"/:id",
		NewRoute(
			handlers.GameItemHandler.FindByID,
			MethodGet,
		),
	)

	gameItemGroup.Add(
		"",
		NewRoute(
			handlers.GameItemHandler.FindAllPaged,
			MethodGet,
		),
	)

	gameItemGroup.Add(
		"/:id",
		NewRoute(
			handlers.GameItemHandler.Update,
			MethodPut,
			WithAccessLevel(access_level.UpdateItem),
		),
	)

	gameItemGroup.Add(
		"/:id",
		NewRoute(
			handlers.GameItemHandler.Delete,
			MethodDelete,
			WithAccessLevel(access_level.DeleteItem),
		),
	)

	return gameItemGroup
}
