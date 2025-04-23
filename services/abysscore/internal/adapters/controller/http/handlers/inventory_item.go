package handlers

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/domain/dto"
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/infrastructure/metrics/tracer"
	"context"
	"github.com/gofiber/fiber/v2"
)

type InventoryItemHandler struct {
	BaseHandler
	inventoryItemService domainservice.InventoryItemService
}

func NewInventoryItemHandler(
	inventoryItemService domainservice.InventoryItemService,
) *InventoryItemHandler {
	return &InventoryItemHandler{inventoryItemService: inventoryItemService}
}

func (h *InventoryItemHandler) GrantInventoryItemToUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := h.extractUser(ctx)
	if err != nil {
		return h.handleError(err, c)
	}

	userID, err := h.extractIntParam("user_id", c)
	if err != nil {
		return h.handleError(err, c)
	}

	itemID, err := h.extractIntParam("item_id", c)
	if err != nil {
		return h.handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(ctx, "inventoryItemService.GrantToUserByAdmin", func(ctx context.Context) (*dto.InventoryItemDTO, error) {
		return h.inventoryItemService.GrantToUserByAdmin(ctx, userID, itemID, admin)
	})
	if err != nil {
		return h.handleError(err, c)
	}

	return h.sendSuccess(result, c)
}

func (h *InventoryItemHandler) GetAllByAuthorization(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := h.extractUser(ctx)
	if err != nil {
		return h.handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(ctx, "inventoryItemService.FindAllByID", func(ctx context.Context) ([]*dto.InventoryItemDTO, error) {
		return h.inventoryItemService.FindAllByUserID(ctx, user.ID)
	})
	if err != nil {
		return h.handleError(err, c)
	}

	return h.sendSuccess(result, c)
}

func (h *InventoryItemHandler) GetAllByUserID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userID, err := h.extractIntParam("user_id", c)
	if err != nil {
		return h.handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(ctx, "inventoryItemService.FindAllByUserID", func(ctx context.Context) ([]*dto.InventoryItemDTO, error) {
		return h.inventoryItemService.FindAllByUserID(ctx, userID)
	})

	if err != nil {
		return h.handleError(err, c)
	}

	return h.sendSuccess(result, c)
}

func (h *InventoryItemHandler) RevokeByAdmin(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := h.extractUser(ctx)
	if err != nil {
		return h.handleError(err, c)
	}

	userID, err := h.extractIntParam("user_id", c)
	if err != nil {
		return h.handleError(err, c)
	}

	itemID, err := h.extractIntParam("item_id", c)
	if err != nil {
		return h.handleError(err, c)
	}

	err = tracer.TraceFn(ctx, "inventoryItemService.RevokeByAdmin", func(ctx context.Context) error {
		return h.inventoryItemService.RevokeByAdmin(ctx, userID, itemID, admin)
	})

	if err != nil {
		return h.handleError(err, c)
	}

	return h.sendNoContent(c)
}

func (h *InventoryItemHandler) SetInventoryItem(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := h.extractUser(ctx)
	if err != nil {
		return h.handleError(err, c)
	}

	r := &request.SetItemAsCurrent{}

	err = h.validateRequest(r, c)
	if err != nil {
		return h.handleError(err, c)
	}

	err = h.inventoryItemService.SetInventoryItemAsCurrent(ctx, user, r.InventoryItemID)
	if err != nil {
		return h.handleError(err, c)
	}

	return h.sendNoContent(c)
}
