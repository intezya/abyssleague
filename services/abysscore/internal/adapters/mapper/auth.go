package mapper

import (
	"abysscore/internal/domain/entity"
	"abysscore/internal/infrastructure/ent"
)

func MapToAuthenticationDataFromEnt(u *ent.User) *entity.AuthenticationData {
	return entity.NewAuthenticationData(
		u.ID,
		u.Username,
		u.Password,
		u.HardwareID,
		u.AccountBlockedUntil,
		u.AccountBlockReason,
	)
}
