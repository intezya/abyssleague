package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/request"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
)

type InventoryItemHandler struct {
	inventoryItemService domainservice.InventoryItemService
}

func NewInventoryItemHandler(
	inventoryItemService domainservice.InventoryItemService,
) *InventoryItemHandler {
	return &InventoryItemHandler{inventoryItemService: inventoryItemService}
}

// GrantInventoryItemToUser grants an item to a user
//
//	@Summary		Grant item to user
//	@Description	Admin grants a game item to a specific user
//	@Tags			Inventory Items
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user_id	path		int											true	"User ID"
//	@Param			item_id	path		int											true	"Item ID"
//	@Success		200		{object}	examples.InventoryItemDTOSuccessResponse	"Granted inventory item"
//	@Failure		400		{object}	examples.BadRequestResponse					"Bad request - invalid ID"
//	@Failure		403		{object}	examples.ForbiddenByAccessLevelResponse		"Forbidden - not enough rights"
//	@Failure		404		{object}	examples.UserNotFoundResponse				"Not found - user not found"
//	@Failure		404		{object}	examples.GameItemNotFound					"Not found - item not found"
//	@Router			/api/users/{user_id}/inventory [post].
func (h *InventoryItemHandler) GrantInventoryItemToUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	userID, err := extractIntParam("user_id", c)
	if err != nil {
		return handleError(err, c)
	}

	itemID, err := extractIntParam("item_id", c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"inventoryItemService.GrantToUserByAdmin",
		func(ctx context.Context) (*dto.InventoryItemDTO, error) {
			return h.inventoryItemService.GrantToUserByAdmin(ctx, userID, itemID, admin)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// GetAllByAuthorization gets all inventory items for the authenticated user
//
//	@Summary		Get current user's inventory
//	@Description	Returns all inventory items for the currently authenticated user
//	@Tags			Inventory Items
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	examples.PaginatedInventoryItemsDTOResponse	"List of user's inventory items"
//	@Router			/api/users/inventory [get].
func (h *InventoryItemHandler) GetAllByAuthorization(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"inventoryItemService.FindAllByID",
		func(ctx context.Context) ([]*dto.InventoryItemDTO, error) {
			return h.inventoryItemService.FindAllByUserID(ctx, user.ID)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// GetAllByUserID gets all inventory items for a specific user
//
//	@Summary		Get user's inventory
//	@Description	Admin retrieves all inventory items for a specified user
//	@Tags			Inventory Items
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user_id	path		int											true	"User ID"
//	@Success		200		{array}		examples.PaginatedInventoryItemsDTOResponse	"List of user's inventory items"
//	@Failure		400		{object}	examples.BadRequestResponse					"Bad request - invalid ID"
//	@Failure		403		{object}	examples.ForbiddenByAccessLevelResponse		"Forbidden - not enough rights"
//	@Failure		404		{object}	examples.UserNotFoundResponse				"Not found - user not found"
//	@Router			/api/users/{user_id}/inventory [get].
func (h *InventoryItemHandler) GetAllByUserID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userID, err := extractIntParam("user_id", c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"inventoryItemService.FindAllByUserID",
		func(ctx context.Context) ([]*dto.InventoryItemDTO, error) {
			return h.inventoryItemService.FindAllByUserID(ctx, userID)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// RevokeByAdmin revokes an item from a user
//
//	@Summary		Revoke item from user
//	@Description	Admin revokes a game item from a specific user
//	@Tags			Inventory Items
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user_id	path	int	true	"User ID"
//	@Param			item_id	path	int	true	"Item ID"
//	@Success		204		"Item successfully revoked"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - invalid ID"
//	@Failure		403		{object}	examples.ForbiddenByAccessLevelResponse	"Forbidden - not enough rights"
//	@Failure		404		{object}	examples.UserNotFoundResponse			"Not found - user not found"
//	@Failure		404		{object}	examples.InventoryItemNotFoundResponse	"Not found - inventory item not found"
//	@Router			/api/users/{user_id}/inventory/{item_id} [delete].
func (h *InventoryItemHandler) RevokeByAdmin(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	userID, err := extractIntParam("user_id", c)
	if err != nil {
		return handleError(err, c)
	}

	itemID, err := extractIntParam("item_id", c)
	if err != nil {
		return handleError(err, c)
	}

	err = tracer.TraceFn(
		ctx,
		"inventoryItemService.RevokeByAdmin",
		func(ctx context.Context) error {
			return h.inventoryItemService.RevokeByAdmin(ctx, userID, itemID, admin)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendNoContent(c)
}

// SetInventoryItem sets an inventory item as current for the user
//
//	@Summary		Set inventory item as current
//	@Description	Sets a specific inventory item as current for the authenticated user
//	@Tags			Inventory Items
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	request.SetItemAsCurrent	true	"Inventory item details"
//	@Success		204		"Item successfully set as current"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - invalid request body"
//	@Failure		404		{object}	examples.InventoryItemNotFoundResponse	"Not found - inventory item not found"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Router			/api/users/me/inventory/set_item [post].
func (h *InventoryItemHandler) SetInventoryItem(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	req, err := getRequest[request.SetItemAsCurrent](c)
	if err != nil {
		return handleError(err, c)
	}

	err = tracer.TraceFn(
		ctx,
		"inventoryItemService.SetInventoryItemAsCurrent",
		func(ctx context.Context) error {
			return h.inventoryItemService.SetInventoryItemAsCurrent(ctx, user, req.InventoryItemID)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendNoContent(c)
}
