package dto

import (
	"time"
)

type InventoryItemDTO struct {
	ID             int       `json:"id"`
	UserID         int       `json:"-"`
	ItemID         int       `json:"-"`
	ReceivedFromID *int      `json:"-"`
	ObtainedAt     time.Time `json:"obtained_at"`

	// Edges
	GameItemID int       `json:"game_item_id"`
	Name       string    `json:"name"`
	Collection string    `json:"collection"`
	Type       int       `json:"type"`
	Rarity     int       `json:"rarity"`
	CreatedAt  time.Time `json:"-"`
}

type GrantInventoryItemDTO struct {
	UserID int `json:"user_id"`
	ItemID int `json:"item_id"`
}

func (g *GrantInventoryItemDTO) ToCreateDTO(performerID int) *CreateInventoryItemDTO {
	return &CreateInventoryItemDTO{
		UserID:         g.UserID,
		ItemID:         g.ItemID,
		ReceivedFromID: performerID,
	}
}

type CreateInventoryItemDTO struct {
	UserID         int
	ItemID         int
	ReceivedFromID int
}
