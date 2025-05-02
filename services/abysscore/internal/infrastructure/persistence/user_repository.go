package persistence

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
	"strings"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/mapper"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	entUser "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/user"
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

// Create adds a new user to the database.
func (r *UserRepository) Create(
	ctx context.Context,
	credentials *dto.CredentialsDTO,
) (*entity.AuthenticationData, error) {
	user, err := r.client.User.
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

// FindDTOById retrieves basic user data by ID.
func (r *UserRepository) FindDTOById(ctx context.Context, id int) (*dto.UserDTO, error) {
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

// FindAuthenticationByLowerUsername retrieves authentication data by lowercase username.
func (r *UserRepository) FindAuthenticationByLowerUsername(
	ctx context.Context,
	lowerUsername string,
) (
	*entity.AuthenticationData,
	error,
) {
	user, err := r.client.User.
		Query().
		Where(entUser.UsernameEqualFold(lowerUsername)).
		Only(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToAuthenticationDataFromEnt(user), nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	exists, _ := r.client.User.
		Query().Where(entUser.EmailEqualFold(email)).
		Exist(ctx)

	return exists
}

// UpdateHWIDByID updates a user's hardware ID.
func (r *UserRepository) UpdateHWIDByID(ctx context.Context, id int, hardwareID string) error {
	_, err := r.client.User.
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

// UpdateLoginStreakLoginAtByID updates a user's login streak and login time.
func (r *UserRepository) UpdateLoginStreakLoginAtByID(
	ctx context.Context,
	id int,
	loginStreak int,
	loginAt time.Time,
) error {
	_, err := r.client.User.
		UpdateOneID(id).
		SetLoginStreak(loginStreak).
		SetLoginAt(loginAt).
		Save(ctx)

	return r.handleUpdateError(err)
}

// UpdatePasswordByID updates a user's password and returns the updated user data.
func (r *UserRepository) UpdatePasswordByID(
	ctx context.Context,
	id int,
	password string,
) (*dto.UserFullDTO, error) {
	user, err := r.client.User.
		UpdateOneID(id).
		SetPassword(password).
		Save(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToUserFullDTOFromEnt(user), nil
}

// SetBlockUntilAndLevelAndReasonFromUser updates a user's block status information.
func (r *UserRepository) SetBlockUntilAndLevelAndReasonFromUser(
	ctx context.Context,
	user *dto.UserDTO,
) error {
	_, err := r.client.User.
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
