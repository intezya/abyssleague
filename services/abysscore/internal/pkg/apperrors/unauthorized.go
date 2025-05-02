package apperrors

import (
	"errors"

	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	errWrongPassword            = errors.New("wrong password")
	errUserWrongHardwareID      = errors.New("wrong hardware id")
	errTokenHardwareIDIsInvalid = errors.New("wrong token hardware id")
)

var (
	ErrWrongPassword            = errorz.Unauthorized(errWrongPassword)
	ErrUserWrongHardwareID      = errorz.Unauthorized(errUserWrongHardwareID)
	ErrTokenHardwareIDIsInvalid = errorz.Unauthorized(errTokenHardwareIDIsInvalid)

	WrapUnauthorized = func(err error) *errorz.Error {
		return errorz.Unauthorized(err)
	}
)
