package mapper

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
)

func MapToUserItemDTOFromEnt(u *ent.UserItem) *dto.UserItemDTO {
	if u == nil {
		return nil
	}

	return &dto.UserItemDTO{
		ID:             u.ID,
		UserID:         u.UserID,
		ItemID:         u.ItemID,
		ReceivedFromID: u.ReceivedFromID,
		ObtainedAt:     u.ObtainedAt,
		Item:           MapToGameItemDTOFromEnt(u.Edges.Item),
	}
}
