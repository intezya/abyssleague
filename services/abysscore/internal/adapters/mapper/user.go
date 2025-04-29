package mapper

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/pkglib/itertools"
)

func ToUserDTOFromEnt(user *ent.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:                     user.ID,
		Username:               user.Username,
		Email:                  user.Email,
		AccessLevel:            user.AccessLevel,
		GenshinUID:             user.GenshinUID,
		HoyolabLogin:           user.HoyolabLogin,
		CurrentMatchID:         user.CurrentMatchID,
		CurrentItemInProfileID: user.CurrentItemInProfileID,
		AvatarURL:              user.AvatarURL,
		InvitesEnabled:         user.InvitesEnabled,
		LoginAt:                user.LoginAt,
		LoginStreak:            user.LoginStreak,
		CreatedAt:              user.CreatedAt,
		SearchBlockedUntil:     user.SearchBlockedUntil,
		SearchBlockReason:      user.SearchBlockReason,
		SearchBlockedLevel:     user.SearchBlockedLevel,
		AccountBlockedUntil:    user.AccountBlockedUntil,
		AccountBlockReason:     user.AccountBlockReason,
		AccountBlockedLevel:    user.AccountBlockedLevel,
	}
}

func ToUserFullDTOFromEnt(user *ent.User) *dto.UserFullDTO {
	mappedFriends := itertools.Map(
		user.Edges.Friends,
		func(friend *ent.User) *dto.UserDTO {
			return ToUserDTOFromEnt(user)
		},
	)

	mappedItems := itertools.Map(user.Edges.Items, ToInventoryItemDTOFromEnt)

	return &dto.UserFullDTO{
		UserDTO: ToUserDTOFromEnt(user),
		// Edges
		Friends:     mappedFriends,
		Items:       mappedItems,
		CurrentItem: ToInventoryItemDTOFromEnt(user.Edges.CurrentItem),
	}
}
