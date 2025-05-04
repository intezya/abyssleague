package mapper

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
)

func ToBannedHardwareIDFromEnt(banned *ent.BannedHardwareID) *dto.BannedHardwareID {
	return &dto.BannedHardwareID{
		ID:        banned.ID,
		CreatedAt: banned.CreatedAt,
		BanReason: banned.BanReason,
	}
}
