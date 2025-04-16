package repositoryports

import (
	"abysscore/internal/domain/entity/userentity"
)

type UserRepository interface {
	Create(credentials *userentity.CredentialsDTO) (*userentity.AuthenticationData, error)

	FindDTOById(id int) (*userentity.UserDTO, error)
	FindFullDTOById(id int) (*userentity.UserFullDTO, error)

	FindAuthenticationByLowerUsername(lowerUsername string) (*userentity.AuthenticationData, error)

	UpdateHWIDByID(id int, hwid string) error
}
