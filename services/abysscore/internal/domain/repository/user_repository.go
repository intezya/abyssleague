package repositoryports

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
	"time"
)

type UserRepository interface {
	WithTx(ctx context.Context) (*ent.Tx, error)
	FindDTOById(ctx context.Context, id int) (*dto.UserDTO, error)
	FindFullDTOById(ctx context.Context, id int) (*dto.UserFullDTO, error)
	ExistsByEmail(ctx context.Context, email string) bool
	SetEmailIfNil(ctx context.Context, userID int, email string) (*dto.UserDTO, error)

	TxCreate(ctx context.Context, tx *ent.Tx, credentials *dto.CredentialsDTO) (*dto.UserDTO, error)
	TxFindDTOById(ctx context.Context, tx *ent.Tx, id int) (*dto.UserDTO, error)
	TxFindFullDTOByLowerUsername(
		ctx context.Context,
		tx *ent.Tx,
		username string,
	) (*dto.UserFullDTO, error)
	TxFindDTOByLowerUsername(ctx context.Context, tx *ent.Tx, username string) (*dto.UserDTO, error)
	TxUpdateLoginStreakLoginAtByID(
		ctx context.Context,
		tx *ent.Tx,
		id int,
		loginStreak int,
		loginAt time.Time,
	) error
	TxSetBlockUntilAndLevelAndReasonFromUser(
		ctx context.Context,
		tx *ent.Tx,
		user *dto.UserDTO,
	) error
}

type AuthenticationRepository interface {
	WithTx(ctx context.Context) (*ent.Tx, error)
	TxUpdateHardwareIDByID(ctx context.Context, tx *ent.Tx, id int, hardwareID string) error
	TxUpdatePasswordByID(ctx context.Context, tx *ent.Tx, id int, password string) error
}

type InventoryRepository interface {
	SetInventoryItemAsCurrent(
		ctx context.Context,
		user *dto.UserDTO,
		item *dto.InventoryItemDTO,
	) error
}

type BannedHardwareIDRepository interface {
	Create(
		ctx context.Context,
		hardwareID string,
		reason optional.String,
	) (*dto.BannedHardwareID, error)
	FindByHardwareID(ctx context.Context, hardwareID string) (*dto.BannedHardwareID, error)
	DeleteByHardwareID(ctx context.Context, hardwareID string) error

	TxFindByHardwareID(
		ctx context.Context,
		tx *ent.Tx,
		hardwareID string,
	) (*dto.BannedHardwareID, error)
}
