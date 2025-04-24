package adaptererror

import (
	"abysscore/internal/common/errors/base"
	"errors"
)

var (
	TooManyRequests     = base.NewError(errors.New("too many requests"), nil, 429)
	InternalServerError = base.NewError(errors.New("internal server error"), nil, 500)
	BadRequest          = base.NewError(errors.New("bad request"), nil, 400)
	BadRequestFunc      = func(wrapped error) *base.Error {
		return base.NewError(errors.New("bad request"), wrapped, 400)
	}

	UnprocessableEntity = func(wrapped error) *base.Error {
		return base.NewError(errors.New("unprocessable entity"), wrapped, 422)
	}

	Unauthorized = func(wrapped error) *base.Error {
		return base.NewError(errors.New("unauthorized"), wrapped, 401)
	}
	InsufficientAccessLevel = base.NewError(errors.New("insufficient access level"), nil, 403)
	UserMatchStateError     = func(wrapped error) *base.Error {
		return base.NewError(errors.New(wrapped.Error()), wrapped, 403)
	}
)
