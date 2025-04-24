package request

import (
	"abysscore/internal/domain/dto"
)

// AuthenticationRequest provide credentials for user registration/login.
type AuthenticationRequest struct {
	Username   string `example:"my_legendary_username"    json:"username"    validate:"required"`
	Password   string `example:"STr0ngP@55w0rD!_"         json:"password"    validate:"required"`
	HardwareID string `example:"QXV0aGVudGljQU1ENjA3NDA0" json:"hardware_id" validate:"required"`
}

func (a *AuthenticationRequest) ToCredentialsDTO() *dto.CredentialsDTO {
	return dto.NewCredentialsDTO(a.Username, a.Password, a.HardwareID)
}

// PasswordChangeRequest provide credentials for password changing.
type PasswordChangeRequest struct {
	Username    string `example:"my_legendary_username"    json:"username"     validate:"required"`
	OldPassword string `example:"STr0ngP@55w0rD!_"         json:"old_password" validate:"required"`
	NewPassword string `example:"QXV0aGVudGljQU1ENjA3NDA0" json:"new_password" validate:"required"`
}

func (a *PasswordChangeRequest) ToDTO() *dto.ChangePasswordDTO {
	return &dto.ChangePasswordDTO{
		Username:    a.Username,
		OldPassword: a.OldPassword,
		NewPassword: a.NewPassword,
	}
}
