package apperrors

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var (
	WrapUnprocessableEntity = func(wrapped error) error {
		return errorz.UnprocessableEntity(wrapped)
	}
)
