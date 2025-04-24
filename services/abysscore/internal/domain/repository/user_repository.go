package repositoryports

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, credentials *dto.CredentialsDTO) (*entity.AuthenticationData, error)

	FindDTOById(ctx context.Context, id int) (*dto.UserDTO, error)
	FindFullDTOById(ctx context.Context, id int) (*dto.UserFullDTO, error)
	FindAuthenticationByLowerUsername(ctx context.Context, lowerUsername string) (*entity.AuthenticationData, error)

	UpdateHWIDByID(ctx context.Context, id int, hwid string) error
	UpdatePasswordByID(ctx context.Context, id int, password string) (*dto.UserFullDTO, error)
	UpdateLoginStreakLoginAtByID(ctx context.Context, id int, loginStreak int, loginAt time.Time) error

	SetBlockUntilAndLevelAndReasonFromUser(ctx context.Context, user *dto.UserDTO) error
	SetInventoryItemAsCurrent(ctx context.Context, user *dto.UserDTO, item *dto.InventoryItemDTO) error
}
