package request

import "abysscore/internal/domain/entity/userentity"

type AuthenticationRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Hwid     string `json:"hwid" validate:"required"`
}

func (a *AuthenticationRequest) ToCredentialsDTO() *userentity.CredentialsDTO {
	return userentity.NewCredentialsDTO(a.Username, a.Password, a.Hwid)
}
