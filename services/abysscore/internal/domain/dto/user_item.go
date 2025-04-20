package dto

import (
	"time"
)

type UserItemDTO struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	ItemID         int       `json:"-"`
	ReceivedFromID int       `json:"received_from_id"`
	ObtainedAt     time.Time `json:"obtained_at"`

	// Edges
	Item *GameItemDTO `json:"item"`
}
