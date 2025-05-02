package dto

import (
	"time"
)

type InventoryItemDTO struct {
	ID             int       `json:"id"`
	UserID         int       `json:"-"`
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

type CreateInventoryItemDTO struct {
	UserID         int
	ItemID         int
	ReceivedFromID int
}
