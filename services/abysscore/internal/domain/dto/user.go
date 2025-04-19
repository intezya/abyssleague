package dto

import (
	"abysscore/internal/domain/entity"
	"abysscore/internal/infrastructure/ent"
	"abysscore/internal/infrastructure/ent/schema/access_level"
	"github.com/intezya/pkglib/itertools"
	"time"
)

type UserDTO struct {
	ID                     int                      `json:"id"`
	Username               string                   `json:"username"`
	LowerUsername          string                   `json:"-"`
	Email                  *string                  `json:"email"`
	AccessLevel            access_level.AccessLevel `json:"-"`
	GenshinUID             *string                  `json:"genshin_uid"`
	HoyolabLogin           *string                  `json:"hoyolab_login"`
	CurrentMatchID         *int                     `json:"-"`
	CurrentItemInProfileID *int                     `json:"-"`
	AvatarURL              *string                  `json:"avatar_url"`
	InvitesEnabled         bool                     `json:"invites_enabled"`
	LoginAt                time.Time                `json:"-"`
	LoginStreak            int                      `json:"login_streak"`
	CreatedAt              time.Time                `json:"created_at"`
	SearchBlockedUntil     *time.Time               `json:"-"`
	AccountBlockedUntil    *time.Time               `json:"-"`
}

type UserFullDTO struct {
	*UserDTO

	//Statistics []*Statistic `json:"statistics"`
	Friends []*UserDTO `json:"friends"`
	//SentFriendRequests []*FriendRequest `json:"sent_friend_requests"`
	//ReceivedFriendRequests []*FriendRequest `json:"received_friend_requests"`
	Items       []*UserItemDTO `json:"items"`
	CurrentItem *UserItemDTO   `json:"current_item"`
	//CurrentMatch *Match      `json:"current_match"`
}

func MapToUserDTOFromEnt(u *ent.User) *UserDTO {
	return &UserDTO{
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
		AccountBlockedUntil:    u.AccountBlockedUntil,
	}
}

func MapToAuthenticationDataFromEnt(u *ent.User) *entity.AuthenticationData {
	return entity.NewAuthenticationData(u.ID, u.Username, u.Password, u.HardwareID)
}

func MapToUserFullDTOFromEnt(u *ent.User) *UserFullDTO {
	mappedFriends := itertools.Map(
		func(v *ent.User) *UserDTO {
			return MapToUserDTOFromEnt(u)
		},
		u.Edges.Friends,
	)

	mappedItems := itertools.Map(
		func(v *ent.UserItem) *UserItemDTO {
			return MapToUserItemDTOFromEnt(v)
		},
		u.Edges.Items,
	)

	return &UserFullDTO{
		UserDTO: MapToUserDTOFromEnt(u),
		// Edges
		Friends:     mappedFriends,
		Items:       mappedItems,
		CurrentItem: MapToUserItemDTOFromEnt(u.Edges.CurrentItem),
	}
}
