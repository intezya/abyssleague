package gameitementity

import (
	"abysscore/internal/domain/entity/types"
	"abysscore/internal/infrastructure/ent"
	"abysscore/internal/infrastructure/ent/gameitem"
)

type OrderBy string

const (
	OrderByCreatedAt  = "created_at"
	OrderByName       = "name"
	OrderByCollection = "collection"
	OrderByType       = "type"
	OrderByRarity     = "rarity"
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
	}

	if orderType == types.OrderAsc {
		return ent.Asc(field)
	}

	return ent.Desc(field)
}
