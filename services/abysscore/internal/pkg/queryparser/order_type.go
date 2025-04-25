package queryparser

import (
	"strings"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/types"
)

func ParseOrderType(input string) (types.OrderType, error) {
	switch strings.ToLower(input) {
	case "", "asc":
		return types.OrderAsc, nil
	case "desc":
		return types.OrderDesc, nil
	default:
		return "", errOrderTypeParseError
	}
}
