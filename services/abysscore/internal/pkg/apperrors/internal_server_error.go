package apperrors

import (
	"errors"

	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	InternalServerError     = errorz.InternalError(nil)
	WrapInternalServerError = func(err error) *errorz.Error { // Used in recover middleware, so it must return *Error
		var typed *errorz.Error
		ok := errors.As(err, &typed)

		if !ok {
			return errorz.InternalError(err)
		}

		return errorz.InternalError(typed.Detail)
	}
)
