package apperrors

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	InternalServerError     = errorz.InternalError(nil)
	WrapInternalServerError = func(err error) *errorz.Error { // Used in recover middleware, so it must return *Error
		return errorz.InternalError(err)
	}
)
