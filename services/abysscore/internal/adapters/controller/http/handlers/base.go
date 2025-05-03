package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/response"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/middleware"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var wrapInvalidRequestBody = func(wrapped error) error {
	return apperrors.WrapUnprocessableEntity(wrapped)
}

// getAndValidateRequest parses and validates the request body into the given generic struct T.
// returns the parsed struct or a validation/parsing error.
func getAndValidateRequest[T interface{}](c *fiber.Ctx) (*T, error) {
	ctx := c.UserContext()
	ctx, span := tracer.StartSpan(ctx, "getAndValidateRequest")
	defer span.End()

	var request T

	err := c.BodyParser(&request)
	if err != nil {
		return nil, wrapInvalidRequestBody(err)
	}

	err = errorz.ValidateJSON(&request)
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
		return nil, apperrors.InternalServerError
	}

	return user, nil
}

// mustExtractUser retrieves the authenticated user from the context.
// can be used if handler has authentication middleware.
// calls panic if user hasn't been authenticated.
func mustExtractUser(ctx context.Context) *dto.UserDTO {
	user, err := extractUser(ctx)

	if err != nil {
		panic(err)
	}

	return user
}

// extractIntParam extracts an integer route parameter by key.
// returns a BadRequest error if the parameter is missing or invalid.
func extractIntParam(key string, c *fiber.Ctx) (int, error) {
	val, err := c.ParamsInt(key)
	if err != nil {
		return 0, apperrors.WrapBadRequest(err)
	}

	return val, nil
}

// handleError maps and sends a consistent error response based on the error type.
func handleError(err error, c *fiber.Ctx) error {
	return apperrors.HandleError(err, c)
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
