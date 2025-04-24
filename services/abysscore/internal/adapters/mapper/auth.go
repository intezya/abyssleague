package mapper

import (
	"abysscore/internal/domain/entity"
	"abysscore/internal/infrastructure/ent"
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
