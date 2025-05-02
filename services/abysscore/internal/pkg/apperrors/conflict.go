package apperrors

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	WrapUserHardwareIDConflict = func(err error) error {
		return errorz.Conflict("user hardware id conflict", err)
	}

	WrapUserAlreadyExists = func(err error) error {
		return errorz.Conflict("user already exists", err)
	}

	ErrAccountAlreadyHasEmail = errorz.Conflict("account already has linked email", nil)

	ErrEmailConflict = errorz.Conflict("someone account already has this email", nil)
)
