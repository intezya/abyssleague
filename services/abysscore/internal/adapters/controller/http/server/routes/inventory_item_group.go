package routes

import (
	"path"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/handlers"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/schema/access_level"
)

func GetInventoryItemGroup(
	handlers *handlers.DependencyProvider,
	provider *DependencyProvider,
) *RouteGroup {
	inventoryItemGroup := NewRouteGroup(path.Join(provider.apiPrefix, "users"))

	inventoryItemGroup.Add(
		"/:user_id/inventory", NewRoute(
			handlers.InventoryItemHandler.GrantInventoryItemToUser,
			MethodPost,
			WithAccessLevel(access_level.GiveItem),
		),
	)

	inventoryItemGroup.Add(
		"/inventory",
		NewRoute(
			handlers.InventoryItemHandler.GetAllByAuthorization,
			MethodGet,
		),
	)

	inventoryItemGroup.Add(
		"/:user_id/inventory",
		NewRoute(
			handlers.InventoryItemHandler.GetAllByUserID,
			MethodGet,
			WithAccessLevel(access_level.ViewInventory),
		),
	)

	inventoryItemGroup.Add(
		"/:user_id/inventory/:item_id",
		NewRoute(
			handlers.InventoryItemHandler.RevokeByAdmin,
			MethodDelete,
			WithAccessLevel(access_level.RevokeItem),
		),
	)

	inventoryItemGroup.Add(
		"/me/inventory/set_item",
		NewRoute(
			handlers.InventoryItemHandler.SetInventoryItem,
			MethodPost,
		),
	)

	return inventoryItemGroup
}
