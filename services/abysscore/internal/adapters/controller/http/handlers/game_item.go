package handlers

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/adapters/controller/http/dto/response"
	adaptererror "abysscore/internal/common/errors/adapter"
	"abysscore/internal/common/errors/base"
	"abysscore/internal/domain/entity/gameitementity"
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/pkg/queryparser"
	"abysscore/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

// GameItemHandler handles HTTP requests for game items
type GameItemHandler struct {
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

// Create handles creating a new game item
func (g *GameItemHandler) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := extractUserFromContext(ctx)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	var r request.CreateUpdateGameItem
	if err := validator.ValidateJSON(&r, c); err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	result, err := g.gameItemService.Create(ctx, &r, user)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.Success(result, c)
}

// FindByID handles retrieving a game item by ID
func (g *GameItemHandler) FindByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := extractIntParamOrErr("id", c)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	result, err := g.gameItemService.FindByID(ctx, id)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.Success(result, c)
}

// FindAllPaged handles retrieving a paginated list of game items
func (g *GameItemHandler) FindAllPaged(c *fiber.Ctx) error {
	ctx := c.UserContext()

	paginationQuery, err := g.paginationQueryFactory(c)
	if err != nil {
		return base.ParseErrorOrInternalResponse(adaptererror.BadRequestFunc(err), c)
	}

	result, err := g.gameItemService.FindAllPaged(ctx, paginationQuery)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.SuccessPagination(result, c)
}

// Update handles updating an existing game item
func (g *GameItemHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()
	user, err := extractUserFromContext(ctx)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	itemID, err := extractIntParamOrErr("id", c)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	var r request.CreateUpdateGameItem
	if err := validator.ValidateJSON(&r, c); err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	err = g.gameItemService.Update(ctx, itemID, &r, user)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.NoContent(c)
}

// Delete handles deleting a game item
func (g *GameItemHandler) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := extractUserFromContext(ctx)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	itemID, err := extractIntParamOrErr("id", c)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	err = g.gameItemService.Delete(ctx, itemID, user)
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.NoContent(c)
}
