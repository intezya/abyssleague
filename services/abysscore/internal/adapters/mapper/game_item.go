package mapper

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
)

func ToGameItemDTOFromEnt(g *ent.GameItem) *dto.GameItemDTO {
	if g == nil {
		return nil
	}

	return &dto.GameItemDTO{
		ID:         g.ID,
		Name:       g.Name,
		Collection: g.Collection,
		Type:       g.Type,
		Rarity:     g.Rarity,
		CreatedAt:  g.CreatedAt,
	}
}
