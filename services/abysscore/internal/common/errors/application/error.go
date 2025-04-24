package applicationerror

import (
	"abysscore/internal/common/errors/base"
	"errors"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrWrongPassword      = base.NewError(errors.New("wrong password"), nil, fiber.StatusUnauthorized)
	ErrUserWrongHwid      = base.NewError(errors.New("wrong hardware id"), nil, fiber.StatusUnauthorized)
	ErrTokenHwidIsInvalid = base.NewError(errors.New("wrong token hardware id"), nil, fiber.StatusUnauthorized)
	ErrAccountIsLocked    = func(reason *string) error {
		var errorMessage string

		if reason != nil {
			errorMessage = *reason
		}

		return base.NewError(errors.New("account locked"), errors.New(errorMessage), fiber.StatusForbidden)
	}
)
