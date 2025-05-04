package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/request"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/gameitementity"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
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

// Create creates a new game item
//
//	@Summary		Create game item
//	@Description	Admin creates a new game item
//	@Tags			Game Items
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.CreateUpdateGameItem				true	"Game item data"
//	@Success		200		{object}	examples.CreateGameItemDTOSuccessResponse	"Created game item"
//	@Failure		400		{object}	examples.BadRequestResponse					"Bad request - missed request fields"
//	@Failure		403		{object}	examples.ForbiddenByAccessLevelResponse		"Forbidden - not enough rights"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse		"Unprocessable entity - invalid request types"
//	@Router			/api/items [post].
func (h *GameItemHandler) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, span := tracer.StartSpan(ctx, "GameItemHandler.Create")
	defer span.End()

	admin := mustExtractUser(ctx)

	req, err := getAndValidateRequest[request.CreateUpdateGameItem](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := h.gameItemService.Create(ctx, req, admin)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// FindByID gets a game item by ID
//
//	@Summary		Get game item by ID
//	@Description	Returns game item by its ID
//	@Tags			Game Items
//	@Produce		json
//	@Param			id	path		int										true	"Game item ID"
//	@Success		200	{object}	examples.FindGameItemDTOSuccessResponse	"Game item"
//	@Failure		400	{object}	examples.BadRequestResponse				"Bad request - invalid ID"
//	@Failure		404	{object}	examples.GameItemNotFound				"Not found - no such game item"
//	@Router			/api/items/{id} [get].
func (h *GameItemHandler) FindByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, span := tracer.StartSpan(ctx, "GameItemHandler.FindByID")
	defer span.End()

	id, err := extractIntParam("id", c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := h.gameItemService.FindByID(ctx, id)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// FindAllPaged returns paginated game items
//
//	@Summary		List game items
//	@Description	Returns a paginated list of game items with sorting
//	@Tags			Game Items
//	@Produce		json
//	@Param			page		query		int										false	"Page number (default: 1)"
//	@Param			size		query		int										false	"Page size (default: 10)"
//	@Param			order_by	query		string									false	"Field to sort by"	Enums(created_at, name, collection, type, rarity)
//	@Param			order_type	query		string									false	"Sort order"		Enums(asc, desc)
//	@Success		200			{object}	examples.PaginatedGameItemsDTOResponse	"Paginated list of game items"
//	@Failure		400			{object}	examples.BadRequestResponse				"Bad request - invalid query params"
//	@Router			/api/items [get].
func (h *GameItemHandler) FindAllPaged(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, span := tracer.StartSpan(ctx, "GameItemHandler.FindAllPaged")
	defer span.End()

	paginationQuery, err := h.getPaginationQuery(c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := h.gameItemService.FindAllPaged(ctx, paginationQuery)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccessPagination(result, c)
}

// Update updates a game item
//
//	@Summary		Update game item
//	@Description	Admin updates a game item by ID
//	@Tags			Game Items
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path	int								true	"Game item ID"
//	@Param			request	body	request.CreateUpdateGameItem	true	"Updated game item data"
//	@Success		204		"Game item updated"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - missed request fields"
//	@Failure		403		{object}	examples.ForbiddenByAccessLevelResponse	"Forbidden - not enough rights"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Router			/api/items/{id} [put].
func (h *GameItemHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, span := tracer.StartSpan(ctx, "GameItemHandler.Update")
	defer span.End()

	admin := mustExtractUser(ctx)

	itemID, err := extractIntParam("id", c)
	if err != nil {
		return handleError(err, c)
	}

	req, err := getAndValidateRequest[request.CreateUpdateGameItem](c)
	if err != nil {
		return handleError(err, c)
	}

	err = h.gameItemService.Update(ctx, itemID, req, admin)
	if err != nil {
		return handleError(err, c)
	}

	return sendNoContent(c)
}

// Delete handles deleting a game item
// Delete deletes a game item
//
//	@Summary		Delete game item
//	@Description	Admin deletes a game item by ID
//	@Tags			Game Items
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Game item ID"
//	@Success		204	"Game item deleted"
//	@Failure		400	{object}	examples.BadRequestResponse				"Bad request - invalid ID"
//	@Failure		403	{object}	examples.ForbiddenByAccessLevelResponse	"Forbidden - not enough rights"
//	@Failure		422	{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Router			/api/items/{id} [delete].
func (h *GameItemHandler) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, span := tracer.StartSpan(ctx, "GameItemHandler.Delete")
	defer span.End()

	admin := mustExtractUser(ctx)

	itemID, err := extractIntParam("id", c)
	if err != nil {
		return handleError(err, c)
	}

	err = h.gameItemService.Delete(ctx, itemID, admin)
	if err != nil {
		return handleError(err, c)
	}

	return sendNoContent(c)
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
		return nil, apperrors.WrapBadRequest(err)
	}

	return paginationQuery, nil
}
