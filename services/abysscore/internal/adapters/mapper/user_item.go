package mapper

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
)

func ToInventoryItemDTOFromEnt(inventoryItem *ent.InventoryItem) *dto.InventoryItemDTO {
	if inventoryItem == nil || inventoryItem.Edges.Item == nil {
		return nil
	}

	gameItem := ToGameItemDTOFromEnt(inventoryItem.Edges.Item)

	return &dto.InventoryItemDTO{
		ID:             inventoryItem.ID,
		UserID:         inventoryItem.UserID,
		ItemID:         inventoryItem.ItemID,
		ReceivedFromID: inventoryItem.ReceivedFromID,
		ObtainedAt:     inventoryItem.ObtainedAt,
		GameItemID:     gameItem.ID,
		Name:           gameItem.Name,
		Collection:     gameItem.Collection,
		Type:           gameItem.Type,
		Rarity:         gameItem.Rarity,
		CreatedAt:      gameItem.CreatedAt,
	}
}
