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
