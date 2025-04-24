package mapper

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
)

func ToGameItemDTOFromEnt(gameItem *ent.GameItem) *dto.GameItemDTO {
	if gameItem == nil {
		return nil
	}

	return &dto.GameItemDTO{
		ID:         gameItem.ID,
		Name:       gameItem.Name,
		Collection: gameItem.Collection,
		Type:       gameItem.Type,
		Rarity:     gameItem.Rarity,
		CreatedAt:  gameItem.CreatedAt,
	}
}
