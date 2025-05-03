package dto

import (
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
	"time"
)

type BannedHardwareID struct {
	ID         int             `json:"id"`
	HardwareID string          `json:"hardware_id"`
	CreatedAt  time.Time       `json:"created_at"`
	BanReason  optional.String `json:"ban_reason"`
}
