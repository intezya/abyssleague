package handlers

import (
	"abysscore/internal/adapters/controller/http/dto/response"
	"abysscore/internal/adapters/controller/http/middleware"
	adaptererror "abysscore/internal/common/errors/adapter"
	"abysscore/internal/common/errors/base"
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/internal/pkg/validator"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct{}

var invalidRequestBodyError = adaptererror.BadRequestFunc(errors.New("invalid request body"))

// ValidateRequest validates the request body and binds it to the provided struct
func (h *BaseHandler) validateRequest(req interface{}, c *fiber.Ctx) error {
	ctx := c.UserContext()

	err := tracer.TraceFn(ctx, "c.BodyParser", func(ctx context.Context) error {
		return c.BodyParser(req)
	})

	if err != nil {
		return base.ParseErrorOrInternalResponse(invalidRequestBodyError, c)
	}

	err = tracer.TraceFn(ctx, "validator.ValidateJSON", func(ctx context.Context) error {
		return validator.ValidateJSON(req)
	})

	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return nil
}

// ExtractUser extracts the user from the context
func (h *BaseHandler) extractUser(ctx context.Context) (*dto.UserDTO, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(*dto.UserDTO)
	if !ok {
		return nil, adaptererror.InternalServerError
	}
	return user, nil
}

// ExtractIntParam extracts an integer parameter from the URL
func (h *BaseHandler) extractIntParam(key string, c *fiber.Ctx) (int, error) {
	val, err := c.ParamsInt(key)
	if err != nil {
		return 0, adaptererror.BadRequestFunc(err)
	}
	return val, nil
}

// HandleError handles errors consistently across all handlers
func (h *BaseHandler) handleError(err error, c *fiber.Ctx) error {
	return base.ParseErrorOrInternalResponse(err, c)
}

// SendSuccess sends a success response
func (h *BaseHandler) sendSuccess(data interface{}, c *fiber.Ctx) error {
	return response.Success(data, c)
}

// SendNoContent sends a no content response
func (h *BaseHandler) sendNoContent(c *fiber.Ctx) error {
	return response.NoContent(c)
}
