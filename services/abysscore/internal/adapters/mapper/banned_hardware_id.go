package mapper

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
)

func ToBannedHardwareIDFromEnt(banned *ent.BannedHardwareID) *dto.BannedHardwareID {
	return &dto.BannedHardwareID{
		ID:         banned.ID,
		HardwareID: banned.HardwareID,
		CreatedAt:  banned.CreatedAt,
		BanReason:  optional.NewP(banned.BanReason),
	}
}
