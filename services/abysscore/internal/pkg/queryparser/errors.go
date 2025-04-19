package queryparser

import (
	"errors"
)

var orderByParseError = errors.New("invalid value for OrderBy")
var orderTypeParseError = errors.New("invalid value for OrderType")
