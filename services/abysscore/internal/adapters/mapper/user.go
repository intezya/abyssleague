package mapper

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
	"github.com/intezya/pkglib/itertools"
)

func MapToUserDTOFromEnt(u *ent.User) *dto.UserDTO {
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

func MapToUserFullDTOFromEnt(u *ent.User) *dto.UserFullDTO {
	mappedFriends := itertools.Map(
		func(v *ent.User) *dto.UserDTO {
			return MapToUserDTOFromEnt(u)
		},
		u.Edges.Friends,
	)

	mappedItems := itertools.Map(
		func(v *ent.UserItem) *dto.UserItemDTO {
			return MapToUserItemDTOFromEnt(v)
		},
		u.Edges.Items,
	)

	return &dto.UserFullDTO{
		UserDTO: MapToUserDTOFromEnt(u),
		// Edges
		Friends:     mappedFriends,
		Items:       mappedItems,
		CurrentItem: MapToUserItemDTOFromEnt(u.Edges.CurrentItem),
	}
}
