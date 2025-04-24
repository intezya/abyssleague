package mapper

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
)

func ToAuthenticationDataFromEnt(user *ent.User) *entity.AuthenticationData {
	return entity.NewAuthenticationData(
		user.ID,
		user.Username,
		user.Password,
		user.HardwareID,
		user.AccountBlockedUntil,
		user.AccountBlockReason,
	)
}
