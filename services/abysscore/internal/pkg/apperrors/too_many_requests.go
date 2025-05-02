package apperrors

import (
	"errors"

	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
)

var errTooManyEmailLinkRequests = errors.New("too many email link requests")

var (
	TooManyEmailLinkRequests = errorz.TooManyRequests(errTooManyEmailLinkRequests)

	TooManyRequests = errorz.TooManyRequests(nil)
)
