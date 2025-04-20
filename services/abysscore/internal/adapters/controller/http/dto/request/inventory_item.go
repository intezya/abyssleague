package request

import "abysscore/internal/domain/dto"

type GrantInventoryItemToUser struct {
	UserID int `json:"user_id"`
	ItemID int `json:"item_id"`
}

func (g *GrantInventoryItemToUser) ToCreateDTO(performerID int) *dto.CreateInventoryItemDTO {
	return &dto.CreateInventoryItemDTO{
		UserID:         g.UserID,
		ItemID:         g.ItemID,
		ReceivedFromID: performerID,
	}
}
