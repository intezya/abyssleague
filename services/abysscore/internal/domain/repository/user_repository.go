package repositoryports

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/userentity"
	"time"
)

type UserRepository interface {
	Create(credentials *userentity.CredentialsDTO) (*userentity.AuthenticationData, error)

	FindDTOById(id int) (*dto.UserDTO, error)
	FindFullDTOById(id int) (*dto.UserFullDTO, error)

	FindAuthenticationByLowerUsername(lowerUsername string) (*userentity.AuthenticationData, error)

	UpdateHWIDByID(id int, hwid string) error

	SetLoginStreakLoginAtByID(id int, loginStreak int, loginAt time.Time) error
}
