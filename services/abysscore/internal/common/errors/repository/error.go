package repositoryerrors

import (
	"abysscore/internal/common/errors/base"
	"errors"
)

var (
	WrapErrUserNotFound = func(err error) *base.Error {
		return base.NewError(errors.New("user not found"), err, 404)
	}
	WrapErrUserHwidConflict = func(err error) *base.Error {
		return base.NewError(errors.New("user hwid conflict"), err, 409)
	}
	WrapErrUserAlreadyExists = func(err error) *base.Error {
		return base.NewError(errors.New("user already exists"), err, 409)
	}

	WrapUnexpectedError = func(err error) *base.Error {
		return base.NewError(errors.New("unexpected error"), err, 500)
	}
)
