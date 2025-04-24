package adaptererror

import (
	"abysscore/internal/common/errors/base"
	"errors"
	"github.com/gofiber/fiber/v2"
)

var (
	TooManyRequests     = base.NewError(errors.New("too many requests"), nil, fiber.StatusTooManyRequests)
	InternalServerError = base.NewError(errors.New("internal server error"), nil, fiber.StatusInternalServerError)
	BadRequest          = base.NewError(errors.New("bad request"), nil, fiber.StatusBadRequest)
	BadRequestFunc      = func(wrapped error) *base.Error {
		return base.NewError(errors.New("bad request"), wrapped, fiber.StatusBadRequest)
	}

	UnprocessableEntity = func(wrapped error) *base.Error {
		return base.NewError(errors.New("unprocessable entity"), wrapped, fiber.StatusUnprocessableEntity)
	}

	Unauthorized = func(wrapped error) *base.Error {
		return base.NewError(errors.New("unauthorized"), wrapped, fiber.StatusUnauthorized)
	}
	InsufficientAccessLevel = base.NewError(errors.New("insufficient access level"), nil, fiber.StatusForbidden)
	UserMatchStateError     = func(wrapped error) *base.Error {
		return base.NewError(errors.New(wrapped.Error()), wrapped, fiber.StatusForbidden)
	}
)
