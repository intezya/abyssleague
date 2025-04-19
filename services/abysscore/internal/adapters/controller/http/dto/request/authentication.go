package request

import (
	"abysscore/internal/domain/entity"
)

type AuthenticationRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Hwid     string `json:"hwid" validate:"required"`
}

func (a *AuthenticationRequest) ToCredentialsDTO() *entity.CredentialsDTO {
	return entity.NewCredentialsDTO(a.Username, a.Password, a.Hwid)
}

type PasswordChangeRequest struct {
	Username    string `json:"username" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (a *PasswordChangeRequest) ToDTO() *entity.ChangePasswordDTO {
	return &entity.ChangePasswordDTO{
		Username:    a.Username,
		OldPassword: a.OldPassword,
		NewPassword: a.NewPassword,
	}
}
