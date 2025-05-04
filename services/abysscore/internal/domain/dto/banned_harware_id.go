package dto

import (
	"time"
)

type BannedHardwareID struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	BanReason *string   `json:"ban_reason"`
}
