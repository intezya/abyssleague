package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/response"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/middleware"
	adaptererror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/adapter"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/base"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/validator"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
)

var wrapInvalidRequestBody = func(wrapped error) error {
	return adaptererror.UnprocessableEntity(wrapped)
}

// getRequest parses and validates the request body into the given generic struct T.
// returns the parsed struct or a validation/parsing error.
func getRequest[T interface{}](c *fiber.Ctx) (*T, error) {
	ctx := c.UserContext()

	var request T

	err := tracer.TraceFn(ctx, "c.BodyParser", func(ctx context.Context) error {
		return c.BodyParser(&request)
	})
	if err != nil {
		return nil, wrapInvalidRequestBody(err)
	}

	err = tracer.TraceFn(ctx, "validator.ValidateJSON", func(ctx context.Context) error {
		return validator.ValidateJSON(&request)
	})
	if err != nil {
		return nil, err
	}

	return &request, nil
}

// extractUser retrieves the authenticated user from the context.
// returns an error if the user is missing or has the wrong type.
func extractUser(ctx context.Context) (*dto.UserDTO, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(*dto.UserDTO)
	if !ok {
		return nil, adaptererror.InternalServerError
	}

	return user, nil
}

// extractIntParam extracts an integer route parameter by key.
// returns a BadRequest error if the parameter is missing or invalid.
func extractIntParam(key string, c *fiber.Ctx) (int, error) {
	val, err := c.ParamsInt(key)
	if err != nil {
		return 0, adaptererror.BadRequestFunc(err)
	}

	return val, nil
}

// handleError maps and sends a consistent error response based on the error type.
func handleError(err error, c *fiber.Ctx) error {
	return base.ParseErrorOrInternalResponse(err, c)
}

// sendSuccess sends a standard JSON success response with the given data.
func sendSuccess(data interface{}, c *fiber.Ctx) error {
	return response.Success(data, c)
}

// sendNoContent sends a 204 No Content response.
func sendNoContent(c *fiber.Ctx) error {
	return response.NoContent(c)
}

// sendSuccessPagination sends a paginated JSON response with the given data.
func sendSuccessPagination[T any](data *dto.PaginatedResult[T], c *fiber.Ctx) error {
	return response.SuccessPagination(data, c)
}
