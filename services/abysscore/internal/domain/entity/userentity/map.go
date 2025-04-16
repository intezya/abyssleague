package userentity

import (
	"abysscore/internal/infrastructure/ent"
)

func MapToDTOFromEnt(u *ent.User) *UserDTO {
	return &UserDTO{
		ID:                     u.ID,
		Username:               u.Username,
		LowerUsername:          u.LowerUsername,
		Email:                  u.Email,
		AccessLevel:            u.AccessLevel,
		GenshinUID:             u.GenshinUID,
		HoyolabLogin:           u.HoyolabLogin,
		CurrentMatchID:         u.CurrentMatchID,
		CurrentItemInProfileID: u.CurrentItemInProfileID,
		AvatarURL:              u.AvatarURL,
		InvitesEnabled:         u.InvitesEnabled,
		LoginAt:                u.LoginAt,
		LoginStreak:            u.LoginStreak,
		CreatedAt:              u.CreatedAt,
		SearchBlockedUntil:     u.SearchBlockedUntil,
		AccountBlockedUntil:    u.AccountBlockedUntil,
	}
}

func MapToAuthenticationDataFromEnt(u *ent.User) *AuthenticationData {
	return NewAuthenticationData(u.ID, u.Username, u.Password, u.HardwareID)
}

func MapToFullDTOFromEnt(u *ent.User) *UserFullDTO {
	return &UserFullDTO{
		UserDTO: *MapToDTOFromEnt(u),
		// TODO: edges
	}
}
