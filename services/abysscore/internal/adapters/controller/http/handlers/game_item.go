package handlers

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/adapters/controller/http/dto/response"
	adaptererror "abysscore/internal/common/errors/adapter"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/internal/pkg/queryparser"
	"context"
	"github.com/gofiber/fiber/v2"
)

// GameItemHandler handles HTTP requests for game items
type GameItemHandler struct {
	BaseHandler

	gameItemService domainservice.GameItemService

	paginationQueryFactory func(c *fiber.Ctx) (*request.PaginationQuery[gameitementity.OrderBy], error)
}

// NewGameItemHandler creates a new game item handler
func NewGameItemHandler(gameItemService domainservice.GameItemService) *GameItemHandler {
	return &GameItemHandler{
		gameItemService: gameItemService,

		paginationQueryFactory: func(c *fiber.Ctx) (
			*request.PaginationQuery[gameitementity.OrderBy],
			error,
		) {
			return request.NewPaginationQuery[gameitementity.OrderBy](c, queryparser.ParseGameEntityOrderBy)
		},
	}
}

// getPaginationQuery gets pagination query parameters from the request
func (g *GameItemHandler) getPaginationQuery(c *fiber.Ctx) (*request.PaginationQuery[gameitementity.OrderBy], error) {
	paginationQuery, err := request.NewPaginationQuery[gameitementity.OrderBy](c, queryparser.ParseGameEntityOrderBy)

	if err != nil {
		return nil, adaptererror.BadRequestFunc(err)
	}

	return paginationQuery, nil
}

// Create handles creating a new game item
func (g *GameItemHandler) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := g.extractUser(ctx)
	if err != nil {
		return g.handleError(err, c)
	}

	r := &request.CreateUpdateGameItem{}
	if err := g.validateRequest(r, c); err != nil {
		return err
	}

	result, err := tracer.TraceFnWithResult(ctx, "gameItemService.Create", func(ctx context.Context) (*dto.GameItemDTO, error) {
		return g.gameItemService.Create(ctx, r, admin)
	})
	if err != nil {
		return g.handleError(err, c)
	}

	return g.sendSuccess(result, c)
}

// FindByID handles retrieving a game item by ID
func (g *GameItemHandler) FindByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := g.extractIntParam("id", c)
	if err != nil {
		return g.handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(ctx, "gameItemService.FindByID", func(ctx context.Context) (*dto.GameItemDTO, error) {
		return g.gameItemService.FindByID(ctx, id)
	})
	if err != nil {
		return g.handleError(err, c)
	}

	return g.sendSuccess(result, c)
}

// FindAllPaged handles retrieving a paginated list of game items
func (g *GameItemHandler) FindAllPaged(c *fiber.Ctx) error {
	ctx := c.UserContext()

	paginationQuery, err := g.getPaginationQuery(c)
	if err != nil {
		return g.handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(ctx, "gameItemService.FindAllPaged", func(ctx context.Context) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
		return g.gameItemService.FindAllPaged(ctx, paginationQuery)
	})
	if err != nil {
		return g.handleError(err, c)
	}

	// SuccessPagination cannot be moved to BaseHandler because the struct methods do not support a generic type
	return response.SuccessPagination(result, c)
}

// Update handles updating an existing game item
func (g *GameItemHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := g.extractUser(ctx)
	if err != nil {
		return g.handleError(err, c)
	}

	itemID, err := g.extractIntParam("id", c)
	if err != nil {
		return g.handleError(err, c)
	}

	r := &request.CreateUpdateGameItem{}
	if err := g.validateRequest(r, c); err != nil {
		return err
	}

	err = tracer.TraceFn(ctx, "gameItemService.Update", func(ctx context.Context) error {
		return g.gameItemService.Update(ctx, itemID, r, admin)
	})
	if err != nil {
		return g.handleError(err, c)
	}

	return g.sendNoContent(c)
}

// Delete handles deleting a game item
func (g *GameItemHandler) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	admin, err := g.extractUser(ctx)
	if err != nil {
		return g.handleError(err, c)
	}

	itemID, err := g.extractIntParam("id", c)
	if err != nil {
		return g.handleError(err, c)
	}

	err = tracer.TraceFn(ctx, "gameItemService.Delete", func(ctx context.Context) error {
		return g.gameItemService.Delete(ctx, itemID, admin)
	})
	if err != nil {
		return g.handleError(err, c)
	}

	return g.sendNoContent(c)
}
