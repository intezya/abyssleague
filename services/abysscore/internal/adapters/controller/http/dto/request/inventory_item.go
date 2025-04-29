package request

import "github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"

type GrantInventoryItemToUser struct {
	UserID int `json:"user_id" validate:"required"`
	ItemID int `json:"item_id" validate:"required"`
}

func (g *GrantInventoryItemToUser) ToCreateDTO(performerID int) *dto.CreateInventoryItemDTO {
	return &dto.CreateInventoryItemDTO{
		UserID:         g.UserID,
		ItemID:         g.ItemID,
		ReceivedFromID: performerID,
	}
}

type SetItemAsCurrent struct {
	InventoryItemID int `json:"inventory_item_id" validate:"required"`
}
