package queryparser

import (
	"abysscore/internal/domain/entity/gameitementity"
)

func ParseGameEntityOrderBy(s string) (gameitementity.OrderBy, error) {
	switch s {
	case "", "created_at":
		return gameitementity.OrderByCreatedAt, nil
	case "name":
		return gameitementity.OrderByName, nil
	case "collection":
		return gameitementity.OrderByCollection, nil
	case "type":
		return gameitementity.OrderByType, nil
	case "rarity":
		return gameitementity.OrderByRarity, nil
	default:
		return "", errOrderByParseError
	}
}
