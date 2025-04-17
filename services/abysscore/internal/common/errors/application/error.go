package applicationerror

import (
	"abysscore/internal/common/errors/base"
	"errors"
)

var (
	ErrWrongPassword      = base.NewError(errors.New("wrong password"), nil, 401)
	ErrUserWrongHwid      = base.NewError(errors.New("wrong hwid"), nil, 401)
	ErrTokenHwidIsInvalid = base.NewError(errors.New("wrong hwid"), nil, 401)
)
