package mapper

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
)

func ToInventoryItemDTOFromEnt(u *ent.InventoryItem) *dto.InventoryItemDTO {
	if u == nil || u.Edges.Item == nil {
		return nil
	}

	gameItem := ToGameItemDTOFromEnt(u.Edges.Item)

	return &dto.InventoryItemDTO{
		ID:             u.ID,
		UserID:         u.UserID,
		ItemID:         u.ItemID,
		ReceivedFromID: u.ReceivedFromID,
		ObtainedAt:     u.ObtainedAt,
		GameItemID:     gameItem.ID,
		Name:           gameItem.Name,
		Collection:     gameItem.Collection,
		Type:           gameItem.Type,
		Rarity:         gameItem.Rarity,
		CreatedAt:      gameItem.CreatedAt,
	}
}
