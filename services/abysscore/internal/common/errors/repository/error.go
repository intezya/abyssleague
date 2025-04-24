package repositoryerrors

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/base"
)

var (
	WrapUserNotFound = func(err error) error {
		return base.NewError(errors.New("user not found"), err, fiber.StatusNotFound)
	}
	WrapUserHwidConflict = func(err error) error {
		return base.NewError(errors.New("user hardware id conflict"), err, fiber.StatusConflict)
	}
	WrapUserAlreadyExists = func(err error) error {
		return base.NewError(errors.New("user already exists"), err, fiber.StatusConflict)
	}

	WrapUnexpectedError = func(err error) error {
		return base.NewError(errors.New("unexpected error"), err, fiber.StatusServiceUnavailable)
	}

	WrapGameItemNotFound = func(err error) error {
		return base.NewError(errors.New("game item not found"), err, fiber.StatusNotFound)
	}
	WrapItemNotFoundInInventory = func(wrapped error) error {
		return base.NewError(errors.New("item not found in inventory"), nil, fiber.StatusNotFound)
	}
)
