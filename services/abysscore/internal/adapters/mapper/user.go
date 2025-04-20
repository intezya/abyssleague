package mapper

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
	"github.com/intezya/pkglib/itertools"
)

func ToUserDTOFromEnt(u *ent.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:                     u.ID,
		Username:               u.Username,
		LowerUsername:          u.LowerUsername,
		Email:                  u.Email,
		AccessLevel:            u.AccessLevel,
		GenshinUID:             u.GenshinUID,
		HoyolabLogin:           u.HoyolabLogin,
		CurrentMatchID:         u.CurrentMatchID,
		CurrentItemInProfileID: u.CurrentItemInProfileID,
		AvatarURL:              u.AvatarURL,
		InvitesEnabled:         u.InvitesEnabled,
		LoginAt:                u.LoginAt,
		LoginStreak:            u.LoginStreak,
		CreatedAt:              u.CreatedAt,
		SearchBlockedUntil:     u.SearchBlockedUntil,
		SearchBlockReason:      u.SearchBlockReason,
		SearchBlockedLevel:     u.SearchBlockedLevel,
		AccountBlockedUntil:    u.AccountBlockedUntil,
		AccountBlockReason:     u.AccountBlockReason,
		AccountBlockedLevel:    u.AccountBlockedLevel,
	}
}

func ToUserFullDTOFromEnt(u *ent.User) *dto.UserFullDTO {
	mappedFriends := itertools.Map(
		u.Edges.Friends,
		func(v *ent.User) *dto.UserDTO {
			return ToUserDTOFromEnt(u)
		},
	)

	mappedItems := itertools.Map(
		u.Edges.Items,
		func(v *ent.InventoryItem) *dto.InventoryItemDTO {
			return ToInventoryItemDTOFromEnt(v)
		},
	)

	return &dto.UserFullDTO{
		UserDTO: ToUserDTOFromEnt(u),
		// Edges
		Friends:     mappedFriends,
		Items:       mappedItems,
		CurrentItem: ToInventoryItemDTOFromEnt(u.Edges.CurrentItem),
	}
}
