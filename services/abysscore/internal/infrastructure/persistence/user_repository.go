package persistence

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"strings"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/mapper"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	entUser "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/user"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
)

// UserRepository provides access to user data in the database.
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new user repository instance.
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// FindDTOById retrieves basic user data by ID.
func (r *UserRepository) FindDTOById(ctx context.Context, id int) (*dto.UserDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.FindDTOById")
	defer span.End()

	user, err := r.client.User.
		Query().
		Where(entUser.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToUserDTOFromEnt(user), nil
}

// FindFullDTOById retrieves complete user data with relationships by ID.
func (r *UserRepository) FindFullDTOById(ctx context.Context, id int) (*dto.UserFullDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.FindFullDTOById")
	defer span.End()

	user, err := r.client.User.
		Query().
		Where(entUser.IDEQ(id)).
		WithCurrentMatch().
		WithFriends().
		WithCurrentItem(func(q *ent.InventoryItemQuery) {
			q.WithItem()
		}).
		WithItems(func(q *ent.InventoryItemQuery) {
			q.WithItem()
		}).
		WithReceivedFriendRequests().
		WithStatistics().
		Only(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToUserFullDTOFromEnt(user), nil
}

// TxFindAuthenticationByLowerUsername retrieves authentication data by lowercase username.
func (r *UserRepository) TxFindAuthenticationByLowerUsername(
	ctx context.Context,
	tx *ent.Tx,
	lowerUsername string,
) (
	*entity.AuthenticationData,
	error,
) {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.TxFindAuthenticationByLowerUsername")
	defer span.End()

	user, err := tx.User.
		Query().
		Where(entUser.UsernameEqualFold(lowerUsername)).
		Only(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToAuthenticationDataFromEnt(user), nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.ExistsByEmail")
	defer span.End()

	exists, _ := r.client.User.
		Query().Where(entUser.EmailEqualFold(email)).
		Exist(ctx)

	return exists
}

// TxUpdateHardwareIDByID updates a user's hardware ID.
func (r *UserRepository) TxUpdateHardwareIDByID(ctx context.Context, tx *ent.Tx, id int, hardwareID string) error {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.TxUpdateHardwareIDByID")
	defer span.End()

	_, err := tx.User.
		UpdateOneID(id).
		SetHardwareID(hardwareID).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) && strings.Contains(err.Error(), "hardwareID") {
			return apperrors.WrapUserHardwareIDConflict(err)
		}

		return r.handleUpdateError(err)
	}

	return nil
}

// TxUpdateLoginStreakLoginAtByID updates a user's login streak and login time.
func (r *UserRepository) TxUpdateLoginStreakLoginAtByID(
	ctx context.Context,
	tx *ent.Tx,
	id int,
	loginStreak int,
	loginAt time.Time,
) error {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.TxUpdateLoginStreakLoginAtByID")
	defer span.End()

	_, err := tx.User.
		UpdateOneID(id).
		SetLoginStreak(loginStreak).
		SetLoginAt(loginAt).
		Save(ctx)

	return r.handleUpdateError(err)
}

// TxUpdatePasswordByID updates a user's password and returns the updated user data.
func (r *UserRepository) TxUpdatePasswordByID(
	ctx context.Context,
	tx *ent.Tx,
	id int,
	password string,
) (*dto.UserFullDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.TxUpdatePasswordByID")
	defer span.End()

	user, err := tx.User.
		UpdateOneID(id).
		SetPassword(password).
		Save(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToUserFullDTOFromEnt(user), nil
}

// TxSetBlockUntilAndLevelAndReasonFromUser updates a user's block status information.
func (r *UserRepository) TxSetBlockUntilAndLevelAndReasonFromUser(
	ctx context.Context,
	tx *ent.Tx,
	user *dto.UserDTO,
) error {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.TxSetBlockUntilAndLevelAndReasonFromUser")
	defer span.End()

	_, err := tx.User.
		UpdateOneID(user.ID).
		SetNillableAccountBlockedUntil(user.AccountBlockedUntil).
		SetAccountBlockedLevel(user.AccountBlockedLevel).
		SetNillableAccountBlockReason(user.AccountBlockReason).
		SetNillableSearchBlockedUntil(user.SearchBlockedUntil).
		SetSearchBlockedLevel(user.SearchBlockedLevel).
		SetNillableSearchBlockReason(user.SearchBlockReason).
		Save(ctx)

	return r.handleUpdateError(err)
}

func (r *UserRepository) SetInventoryItemAsCurrent(
	ctx context.Context,
	user *dto.UserDTO,
	item *dto.InventoryItemDTO,
) error {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.SetInventoryItemAsCurrent")
	defer span.End()

	_, err := r.client.User.
		UpdateOneID(user.ID).
		SetCurrentItemID(item.ID).
		Save(ctx)
	if err != nil {
		return r.handleUpdateError(err)
	}

	return nil
}

func (r *UserRepository) SetEmailIfNil(
	ctx context.Context,
	userID int,
	email string,
) (*dto.UserDTO, error) {
	return withTxResult(
		ctx, r.client, func(tx *ent.Tx) (*dto.UserDTO, error) {
			user, err := tx.User.Get(ctx, userID)
			if err != nil {
				return nil, r.handleQueryError(err)
			}

			if user.Email != nil {
				return nil, apperrors.ErrAccountAlreadyHasEmail
			}

			savedUser, err := tx.User.UpdateOneID(userID).SetEmail(email).Save(ctx)
			if err != nil {
				return nil, r.handleConstraintError(err) // unexpected
			}

			return mapper.ToUserDTOFromEnt(savedUser), nil
		},
	)
}

func (r *UserRepository) WithTx(ctx context.Context) (*ent.Tx, error) {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.WithTx")
	defer span.End()

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, apperrors.WrapUnexpectedError(err)
	}

	return tx, nil
}

// TxCreate adds a new user to the database.
func (r *UserRepository) TxCreate(
	ctx context.Context,
	tx *ent.Tx,
	credentials *dto.CredentialsDTO,
) (*entity.AuthenticationData, error) {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.TxCreate")
	defer span.End()

	user, err := tx.User.
		Create().
		SetUsername(credentials.Username).
		SetPassword(credentials.Password).
		SetHardwareID(credentials.HardwareID).
		Save(ctx)
	if err != nil {
		return nil, r.handleConstraintError(err)
	}

	return mapper.ToAuthenticationDataFromEnt(user), nil
}

// TxFindDTOById retrieves basic user data by ID.
func (r *UserRepository) TxFindDTOById(ctx context.Context, tx *ent.Tx, id int) (*dto.UserDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "UserRepository.TxFindDTOById")
	defer span.End()

	user, err := tx.User.
		Query().
		Where(entUser.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToUserDTOFromEnt(user), nil
}

// Helper methods for error handling

// handleQueryError transforms Ent query errors into domain-specific errors.
func (r *UserRepository) handleQueryError(err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return apperrors.WrapUserNotFound(err)
	}

	return apperrors.WrapUnexpectedError(err)
}

// handleUpdateError transforms Ent update errors into domain-specific errors.
func (r *UserRepository) handleUpdateError(err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return apperrors.WrapUserNotFound(err)
	}

	return apperrors.WrapUnexpectedError(err)
}

// handleConstraintError processes database constraint violation errors.
func (r *UserRepository) handleConstraintError(err error) error {
	if err == nil {
		return nil
	}

	if !ent.IsConstraintError(err) {
		return apperrors.WrapUnexpectedError(err)
	}

	switch {
	case strings.Contains(err.Error(), "username"):
		return apperrors.WrapUserAlreadyExists(err)
	case strings.Contains(err.Error(), "hardware_id"):
		return apperrors.WrapUserHardwareIDConflict(err)
	case strings.Contains(err.Error(), "email"):
		return apperrors.ErrAccountAlreadyHasEmail
	default:
		return apperrors.WrapUnexpectedError(err)
	}
}
