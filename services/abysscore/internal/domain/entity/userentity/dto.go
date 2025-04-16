package userentity

import (
	"abysscore/internal/infrastructure/ent/schema/access_level"
	"time"
)

type UserDTO struct {
	ID                     int                      `json:"id,omitempty"`
	Username               string                   `json:"username,omitempty"`
	LowerUsername          string                   `json:"lower_username,omitempty"`
	Email                  *string                  `json:"email,omitempty"`
	AccessLevel            access_level.AccessLevel `json:"access_level,omitempty"`
	GenshinUID             *string                  `json:"genshin_uid,omitempty"`
	HoyolabLogin           *string                  `json:"hoyolab_login,omitempty"`
	CurrentMatchID         *int                     `json:"current_match_id,omitempty"`
	CurrentItemInProfileID *int                     `json:"current_item_in_profile_id,omitempty"`
	AvatarURL              *string                  `json:"avatar_url,omitempty"`
	InvitesEnabled         bool                     `json:"invites_enabled,omitempty"`
	LoginAt                time.Time                `json:"login_at,omitempty"`
	LoginStreak            int                      `json:"login_streak,omitempty"`
	CreatedAt              time.Time                `json:"created_at,omitempty"`
	SearchBlockedUntil     *time.Time               `json:"search_blocked_until,omitempty"`
	AccountBlockedUntil    *time.Time               `json:"account_blocked_until,omitempty"`
	//TODO: Edges        UserEdges `json:"edges"`
}

type UserFullDTO struct {
	UserDTO
	// TODO: Edges
}
