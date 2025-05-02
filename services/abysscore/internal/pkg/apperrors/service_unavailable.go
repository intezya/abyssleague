package apperrors

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var WrapServiceUnavailable = func(err error) error {
	return errorz.ServiceUnavailable(err)
}
