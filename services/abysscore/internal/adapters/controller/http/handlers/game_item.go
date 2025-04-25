package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/request"
	adaptererror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/adapter"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/gameitementity"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/queryparser"
)

// GameItemHandler handles HTTP requests for game items.
type GameItemHandler struct {
	gameItemService domainservice.GameItemService

	paginationQueryFactory func(c *fiber.Ctx) (*request.PaginationQuery[gameitementity.OrderBy], error)
}

// NewGameItemHandler creates a new game item handler.
func NewGameItemHandler(gameItemService domainservice.GameItemService) *GameItemHandler {
	return &GameItemHandler{
		gameItemService: gameItemService,

		paginationQueryFactory: func(c *fiber.Ctx) (
			*request.PaginationQuery[gameitementity.OrderBy],
			error,
		) {
			return request.NewPaginationQuery[gameitementity.OrderBy](
				c,
				queryparser.ParseGameEntityOrderBy,
			)
		},
	}
}

// getPaginationQuery gets pagination query parameters from the request.
func (h *GameItemHandler) getPaginationQuery(
	c *fiber.Ctx,
) (*request.PaginationQuery[gameitementity.OrderBy], error) {
	paginationQuery, err := request.NewPaginationQuery[gameitementity.OrderBy](
		c,
		queryparser.ParseGameEntityOrderBy,
	)
	if err != nil {
		return nil, adaptererror.BadRequestFunc(err)
	}

	return paginationQuery, nil
}

// @Router /api/items [post].
func (h *GameItemHandler) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	req, err := getRequest[request.CreateUpdateGameItem](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"gameItemService.Create",
		func(ctx context.Context) (*dto.GameItemDTO, error) {
			return h.gameItemService.Create(ctx, req, admin)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// @Router /api/items/{id} [get].
func (h *GameItemHandler) FindByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := extractIntParam("id", c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"gameItemService.FindByID",
		func(ctx context.Context) (*dto.GameItemDTO, error) {
			return h.gameItemService.FindByID(ctx, id)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// @Router /api/items [get].
func (h *GameItemHandler) FindAllPaged(c *fiber.Ctx) error {
	ctx := c.UserContext()

	paginationQuery, err := h.getPaginationQuery(c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"gameItemService.FindAllPaged",
		func(ctx context.Context) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
			return h.gameItemService.FindAllPaged(ctx, paginationQuery)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccessPagination(result, c)
}

// @Router /api/items/{id} [put].
func (h *GameItemHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	itemID, err := extractIntParam("id", c)
	if err != nil {
		return handleError(err, c)
	}

	req, err := getRequest[request.CreateUpdateGameItem](c)
	if err != nil {
		return handleError(err, c)
	}

	err = tracer.TraceFn(
		ctx,
		"gameItemService.Update",
		func(ctx context.Context) error {
			return h.gameItemService.Update(ctx, itemID, req, admin)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendNoContent(c)
}

// @Router /api/items/{id} [delete].
func (h *GameItemHandler) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	itemID, err := extractIntParam("id", c)
	if err != nil {
		return handleError(err, c)
	}

	err = tracer.TraceFn(
		ctx,
		"gameItemService.Delete",
		func(ctx context.Context) error {
			return h.gameItemService.Delete(ctx, itemID, admin)
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendNoContent(c)
}
