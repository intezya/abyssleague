package adaptererror

import (
	"abysscore/internal/common/errors/base"
	"errors"
)

var (
	TooManyRequests     = base.NewError(errors.New("too many requests"), nil, 429)
	InternalServerError = base.NewError(errors.New("internal server error"), nil, 500)
	BadRequest          = base.NewError(errors.New("bad request"), nil, 400)
	ErrUnauthorized     = func(wrapped error) *base.Error {
		return base.NewError(errors.New("unauthorized"), wrapped, 401)
	}
)
