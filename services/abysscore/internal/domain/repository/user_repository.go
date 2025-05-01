package repositoryports

import (
	"context"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
)

type UserRepository interface {
	FindDTOById(ctx context.Context, id int) (*dto.UserDTO, error)
	FindFullDTOById(ctx context.Context, id int) (*dto.UserFullDTO, error)
	ExistsByEmail(ctx context.Context, email string) bool
	SetEmailIfNil(ctx context.Context, userID int, email string) (*dto.UserDTO, error)
}

type AuthenticationRepository interface {
	Create(ctx context.Context, credentials *dto.CredentialsDTO) (*entity.AuthenticationData, error)
	FindAuthenticationByLowerUsername(
		ctx context.Context,
		lowerUsername string,
	) (*entity.AuthenticationData, error)
	UpdatePasswordByID(ctx context.Context, id int, password string) (*dto.UserFullDTO, error)
	UpdateLoginStreakLoginAtByID(
		ctx context.Context,
		id int,
		loginStreak int,
		loginAt time.Time,
	) error
	UpdateHWIDByID(ctx context.Context, id int, hwid string) error
	SetBlockUntilAndLevelAndReasonFromUser(ctx context.Context, user *dto.UserDTO) error
}

type InventoryRepository interface {
	SetInventoryItemAsCurrent(
		ctx context.Context,
		user *dto.UserDTO,
		item *dto.InventoryItemDTO,
	) error
}
