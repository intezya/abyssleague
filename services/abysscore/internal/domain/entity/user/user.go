package user

import (
	"abysscore/internal/domain/entity/user/access_level"
	"time"
)

type User struct {
	ID int `db:"id"`

	Username    string                   `db:"username"`
	Password    string                   `db:"password"`
	HardwareID  string                   `db:"hardware_id"`
	AccessLevel access_level.AccessLevel `db:"access_level"`

	AvatarUrl string `db:"avatar_url"`

	CurrentMatchID int `db:"current_match_id"`

	LoginAt     time.Time `db:"login_at"`
	LoginStreak int       `db:"login_streak"`

	CreatedAt           time.Time  `db:"created_at"`
	SearchBlockedUntil  *time.Time `db:"search_blocked_until"`
	AccountBlockedUntil *time.Time `db:"account_blocked_until"`
}

func (u *User) GetUsername() (username string) {
	return u.Username
}

func (u *User) GetHardwareID() (hwid string) {
	return u.HardwareID
}
