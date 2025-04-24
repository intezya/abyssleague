package request

import (
	"abysscore/internal/domain/dto"
)

// AuthenticationRequest provide credentials for user registration/login.
type AuthenticationRequest struct {
	Username   string `json:"username"    validate:"required" example:"my_legendary_username"`
	Password   string `json:"password"    validate:"required" example:"STr0ngP@55w0rD!_"`
	HardwareID string `json:"hardware_id" validate:"required" example:"QXV0aGVudGljQU1ENjA3NDA0"`
}

func (a *AuthenticationRequest) ToCredentialsDTO() *dto.CredentialsDTO {
	return dto.NewCredentialsDTO(a.Username, a.Password, a.HardwareID)
}

// PasswordChangeRequest provide credentials for password changing.
type PasswordChangeRequest struct {
	Username    string `json:"username"     validate:"required" example:"my_legendary_username"`
	OldPassword string `json:"old_password" validate:"required" example:"STr0ngP@55w0rD!_"`
	NewPassword string `json:"new_password" validate:"required" example:"QXV0aGVudGljQU1ENjA3NDA0"`
}

func (a *PasswordChangeRequest) ToDTO() *dto.ChangePasswordDTO {
	return &dto.ChangePasswordDTO{
		Username:    a.Username,
		OldPassword: a.OldPassword,
		NewPassword: a.NewPassword,
	}
}
