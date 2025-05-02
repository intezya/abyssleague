package apperrors

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	WrapUnexpectedError = func(err error) error {
		return errorz.InternalError(err)
	}
)
