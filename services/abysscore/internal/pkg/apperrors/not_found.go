package apperrors

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	WrapUserNotFound = func(err error) error {
		return errorz.NotFound("user", err)
	}

	WrapGameItemNotFound = func(err error) error {
		return errorz.NotFound("game item", err)
	}

	WrapItemNotFoundInInventory = func(err error) error {
		return errorz.NotFound("item", err)
	}

	WrapMailDataNotFound = func(err error) error {
		return errorz.NotFound("mail data", err)
	}
)
