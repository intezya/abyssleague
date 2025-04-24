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
