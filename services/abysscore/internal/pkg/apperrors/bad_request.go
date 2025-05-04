package apperrors

import (
	"errors"

	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var errWrongVerificationCode = errors.New("wrong verification code")

var (
	ErrWrongVerificationCodeForEmailLink  = errorz.BadRequest(errWrongVerificationCode)
	WrapWrongVerificationCodeForEmailLink = func(err error) error {
		return errorz.ServiceUnavailable(err)
	}

	WrapBadRequest = func(err error) error {
		return errorz.BadRequest(err)
	}
)
