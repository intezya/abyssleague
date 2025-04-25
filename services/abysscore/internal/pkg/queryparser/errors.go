package queryparser

import (
	"errors"
)

var (
	errOrderByParseError   = errors.New("invalid value for OrderBy")
	errOrderTypeParseError = errors.New("invalid value for OrderType")
)
