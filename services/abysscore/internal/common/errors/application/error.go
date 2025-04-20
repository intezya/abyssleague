package applicationerror

import (
	"abysscore/internal/common/errors/base"
	"errors"
)

var (
	ErrWrongPassword      = base.NewError(errors.New("wrong password"), nil, 401)
	ErrUserWrongHwid      = base.NewError(errors.New("wrong hwid"), nil, 401)
	ErrTokenHwidIsInvalid = base.NewError(errors.New("wrong hwid"), nil, 401)
	ErrAccountIsLocked    = func(reason *string) error {
		var r string

		if reason != nil {
			r = *reason
		}

		return base.NewError(errors.New("account locked"), errors.New(r), 403)
	}
	ErrItemNotFoundInInventory = base.NewError(errors.New("item not found in inventory"), nil, 404)
)
