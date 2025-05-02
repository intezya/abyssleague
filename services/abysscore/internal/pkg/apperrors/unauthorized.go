package apperrors

import (
	"errors"
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	ErrWrongPassword            = errorz.Unauthorized(errors.New("wrong password"))
	ErrUserWrongHardwareID      = errorz.Unauthorized(errors.New("wrong hardware id"))
	ErrTokenHardwareIDIsInvalid = errorz.Unauthorized(errors.New("wrong token hardware id"))

	WrapUnauthorized = func(err error) error {
		return errorz.Unauthorized(err)
	}
)
