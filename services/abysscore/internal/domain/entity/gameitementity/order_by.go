package gameitementity

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/types"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/gameitem"
)

type OrderBy string

const (
	OrderByCreatedAt  OrderBy = "created_at"
	OrderByName       OrderBy = "name"
	OrderByCollection OrderBy = "collection"
	OrderByType       OrderBy = "type"
	OrderByRarity     OrderBy = "rarity"
)

func (o OrderBy) ToOrderOption(orderType types.OrderType) gameitem.OrderOption {
	var field string

	switch o {
	case OrderByCreatedAt:
		field = gameitem.FieldCreatedAt
	case OrderByName:
		field = gameitem.FieldName
	case OrderByCollection:
		field = gameitem.FieldCollection
	case OrderByType:
		field = gameitem.FieldType
	case OrderByRarity:
		field = gameitem.FieldRarity
	default:
		field = gameitem.FieldCreatedAt
	}

	if orderType == types.OrderAsc {
		return ent.Asc(field)
	}

	return ent.Desc(field)
}
