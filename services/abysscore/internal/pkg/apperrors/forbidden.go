package apperrors

import (
	"errors"

	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	ErrAccountIsLocked = func(reason *string) error {
		//nolint:err113 // required dynamic error
		var reasonAsError error

		if reason != nil {
			reasonAsError = errors.New(*reason)
		}

		return errorz.Forbidden("account is locked", reasonAsError)
	}

	ErrHardwareIDBanned = func(reason *string) error {
		//nolint:err113 // required dynamic error
		var reasonAsError error

		if reason != nil {
			reasonAsError = errors.New(*reason)
		}

		return errorz.Forbidden("hardware id is banned", reasonAsError)
	}

	ForbiddenByInsufficientAccessLevel = errorz.Forbidden("insufficient access level", nil)

	WrapUserMatchStateError = func(err error) error {
		return errorz.Forbidden("account is locked", err)
	}
)
