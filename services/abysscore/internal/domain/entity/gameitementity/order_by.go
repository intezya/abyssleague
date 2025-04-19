package gameitementity

import "abysscore/internal/infrastructure/ent/gameitem"

type OrderBy string

const (
	OrderByUnset      OrderBy = "unset"
	OrderByCreatedAt          = "created_at"
	OrderByName               = "name"
	OrderByCollection         = "collection"
	OrderByType               = "type"
	OrderByRarity             = "rarity"
)

func (o OrderBy) ToOrderOption() gameitem.OrderOption {
	switch o {
	case OrderByUnset:
		return nil
	case OrderByCreatedAt:
		return gameitem.ByCreatedAt()
	case OrderByName:
		return gameitem.ByName()
	case OrderByCollection:
		return gameitem.ByCollection()
	case OrderByType:
		return gameitem.ByType()
	case OrderByRarity:
		return gameitem.ByRarity()
	default:
		return nil
	}
}
