package repositoryports

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity"
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, credentials *entity.CredentialsDTO) (*entity.AuthenticationData, error)

	FindDTOById(ctx context.Context, id int) (*dto.UserDTO, error)
	FindFullDTOById(ctx context.Context, id int) (*dto.UserFullDTO, error)
	FindAuthenticationByLowerUsername(ctx context.Context, lowerUsername string) (*entity.AuthenticationData, error)

	UpdateHWIDByID(ctx context.Context, id int, hwid string) error
	UpdatePasswordByID(ctx context.Context, id int, password string) (*dto.UserFullDTO, error)
	UpdateLoginStreakLoginAtByID(ctx context.Context, id int, loginStreak int, loginAt time.Time) error

	SetBlockUntilAndLevelAndReasonFromUser(ctx context.Context, user *dto.UserDTO) error
}
