package mapper

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/pkglib/logger"
)

func ToInventoryItemDTOFromEnt(inventoryItem *ent.InventoryItem) *dto.InventoryItemDTO {
	if inventoryItem == nil {
		return nil
	}

	if inventoryItem.Edges.Item == nil {
		logger.Log.Warn("inventory item base is nil, cannot map")
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
