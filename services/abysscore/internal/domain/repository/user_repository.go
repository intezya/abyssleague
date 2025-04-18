package repositoryports

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity"
	"time"
)

type UserRepository interface {
	Create(credentials *entity.CredentialsDTO) (*entity.AuthenticationData, error)

	FindDTOById(id int) (*dto.UserDTO, error)
	FindFullDTOById(id int) (*dto.UserFullDTO, error)

	FindAuthenticationByLowerUsername(lowerUsername string) (*entity.AuthenticationData, error)

	UpdateHWIDByID(id int, hwid string) error

	SetLoginStreakLoginAtByID(id int, loginStreak int, loginAt time.Time) error
}
