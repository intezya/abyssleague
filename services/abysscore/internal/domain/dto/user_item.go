package dto

import (
	"abysscore/internal/infrastructure/ent"
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

func MapToUserItemDTOFromEnt(u *ent.UserItem) *UserItemDTO {
	if u == nil {
		return nil
	}

	return &UserItemDTO{
		ID:             u.ID,
		UserID:         u.UserID,
		ItemID:         u.ItemID,
		ReceivedFromID: u.ReceivedFromID,
		ObtainedAt:     u.ObtainedAt,
		Item:           MapToGameItemDTOFromEnt(u.Edges.Item),
	}
}
