package apperrors

import (
	"errors"
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
)

var (
	ErrAccountIsLocked = func(reason optional.String) error {
		reasonAsError := errors.New(reason.Default("reason is null"))

		return errorz.Forbidden("account is locked", reasonAsError)
	}

	ForbiddenByInsufficientAccessLevel = errorz.Forbidden("insufficient access level", nil)

	WrapUserMatchStateError = func(err error) error {
		return errorz.Forbidden("account is locked", err)
	}
)
