package repositoryerrors

import (
	"abysscore/internal/common/errors/base"
	"errors"
)

var (
	WrapUserNotFound = func(err error) *base.Error {
		return base.NewError(errors.New("user not found"), err, 404)
	}
	WrapUserHwidConflict = func(err error) *base.Error {
		return base.NewError(errors.New("user hwid conflict"), err, 409)
	}
	WrapUserAlreadyExists = func(err error) *base.Error {
		return base.NewError(errors.New("user already exists"), err, 409)
	}

	WrapUnexpectedError = func(err error) *base.Error {
		return base.NewError(errors.New("unexpected error"), err, 502)
	}

	WrapGameItemNotFound = func(err error) *base.Error {
		return base.NewError(errors.New("game item not found"), err, 404)
	}
	WrapItemNotFoundInInventory = func(wrapped error) error {
		return base.NewError(errors.New("item not found in inventory"), nil, 404)
	}
)
