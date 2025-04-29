package dto

import (
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/schema/access_level"
)

type UserDTO struct {
	ID                     int                      `json:"id"`
	Username               string                   `json:"username"`
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

	SearchBlockedUntil *time.Time `json:"-"`
	SearchBlockReason  *string    `json:"-"`
	SearchBlockedLevel int        `json:"-"`

	AccountBlockedUntil *time.Time `json:"-"`
	AccountBlockReason  *string    `json:"-"`
	AccountBlockedLevel int        `json:"-"`
}

type UserFullDTO struct {
	*UserDTO

	// Statistics []*Statistic `json:"statistics"`
	Friends []*UserDTO `json:"friends"`
	// SentFriendRequests []*FriendRequest `json:"sent_friend_requests"`
	// ReceivedFriendRequests []*FriendRequest `json:"received_friend_requests"`
	Items       []*InventoryItemDTO `json:"items"`
	CurrentItem *InventoryItemDTO   `json:"current_item"`
	// CurrentMatch *Match      `json:"current_match"`
}
