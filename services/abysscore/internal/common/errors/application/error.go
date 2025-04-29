package applicationerror

import (
	"errors"
	repositoryerrors "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/base"
)

var (
	ErrWrongPassword = base.NewError(
		errors.New("wrong password"),
		nil,
		fiber.StatusUnauthorized,
	)
	ErrUserWrongHwid = base.NewError(
		errors.New("wrong hardware id"),
		nil,
		fiber.StatusUnauthorized,
	)
	ErrTokenHwidIsInvalid = base.NewError(
		errors.New("wrong token hardware id"),
		nil,
		fiber.StatusUnauthorized,
	)
	ErrAccountIsLocked = func(reason *string) error {
		var errorMessage string

		if reason != nil {
			errorMessage = *reason
		}

		return base.NewError(
			errors.New("account locked"),
			errors.New(errorMessage),
			fiber.StatusForbidden,
		)
	}
	ErrAccountAlreadyHasEmail = repositoryerrors.ErrAccountAlreadyHasEmail
	ErrEmailConflict          = base.NewError(
		errors.New("email conflict"),
		nil,
		fiber.StatusConflict,
	)
	ErrWrongVerificationCodeForEmailLink = base.NewError(
		errors.New("wrong verification code"),
		nil,
		fiber.StatusBadRequest,
	)
	WrapWrongVerificationCodeForEmailLink = func(err error) *base.Error {
		return base.NewError(
			errors.New("wrong verification code"),
			err,
			fiber.StatusBadRequest,
		)
	}
	TooManyEmailLinkRequests = base.NewError(
		errors.New("too many email link requests"),
		nil,
		fiber.StatusTooManyRequests,
	)

	WrapServiceUnavailable = func(err error) *base.Error {
		return base.NewError(errors.New("service unavailable"), err, fiber.StatusServiceUnavailable)
	}
)
