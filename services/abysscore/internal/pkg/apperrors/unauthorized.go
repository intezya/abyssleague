package apperrors

import (
	"errors"

	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	errWrongPassword            = errors.New("wrong password")
	errUserWrongHardwareID      = errors.New("wrong hardware id")
	errTokenHardwareIDIsInvalid = errors.New("wrong token hardware id")
	errHardwareIDIsInvalid      = errors.New("hardware id is invalid. contact the support")
)

var (
	ErrWrongPassword            = errorz.Unauthorized(errWrongPassword)
	ErrUserWrongHardwareID      = errorz.Unauthorized(errUserWrongHardwareID)
	ErrTokenHardwareIDIsInvalid = errorz.Unauthorized(errTokenHardwareIDIsInvalid)
	ErrHardwareIDIsInvalid      = errorz.Unauthorized(errHardwareIDIsInvalid)

	WrapUnauthorized = func(err error) *errorz.Error {
		var typed *errorz.Error
		ok := errors.As(err, &typed)

		if !ok {
			return errorz.Unauthorized(err)
		}

		return errorz.Unauthorized(typed.Detail)
	}
)
