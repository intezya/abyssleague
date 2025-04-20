package persistence

import (
	"abysscore/internal/adapters/mapper"
	repositoryerrors "abysscore/internal/common/errors/repository"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity"
	"abysscore/internal/infrastructure/ent"
	"abysscore/internal/infrastructure/ent/user"
	"context"
	"strings"
	"time"
)

// UserRepository provides access to user data in the database
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// Create adds a new user to the database
func (r *UserRepository) Create(ctx context.Context, credentials *dto.CredentialsDTO) (*entity.AuthenticationData, error) {
	u, err := r.client.User.
		Create().
		SetUsername(credentials.Username).
		SetLowerUsername(strings.ToLower(credentials.Username)).
		SetPassword(credentials.Password).
		SetHardwareID(credentials.Hwid).
		Save(ctx)

	if err != nil {
		return nil, r.handleConstraintError(err)
	}

	return mapper.MapToAuthenticationDataFromEnt(u), nil
}

// FindDTOById retrieves basic user data by ID
func (r *UserRepository) FindDTOById(ctx context.Context, id int) (*dto.UserDTO, error) {
	u, err := r.client.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.MapToUserDTOFromEnt(u), nil
}

// FindFullDTOById retrieves complete user data with relationships by ID
func (r *UserRepository) FindFullDTOById(ctx context.Context, id int) (*dto.UserFullDTO, error) {
	u, err := r.client.User.
		Query().
		Where(user.IDEQ(id)).
		WithCurrentMatch().
		WithFriends().
		WithCurrentItem().
		WithItems().
		WithReceivedFriendRequests().
		WithStatistics().
		Only(ctx)

	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.MapToUserFullDTOFromEnt(u), nil
}

// FindAuthenticationByLowerUsername retrieves authentication data by lowercase username
func (r *UserRepository) FindAuthenticationByLowerUsername(ctx context.Context, lowerUsername string) (
	*entity.AuthenticationData,
	error,
) {
	u, err := r.client.User.
		Query().
		Where(user.LowerUsernameEQ(lowerUsername)).
		Only(ctx)

	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.MapToAuthenticationDataFromEnt(u), nil
}

// UpdateHWIDByID updates a user's hardware ID
func (r *UserRepository) UpdateHWIDByID(ctx context.Context, id int, hwid string) error {
	_, err := r.client.User.
		UpdateOneID(id).
		SetHardwareID(hwid).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) && strings.Contains(err.Error(), "hwid") {
			return repositoryerrors.WrapErrUserHwidConflict(err)
		}
		return r.handleUpdateError(err)
	}

	return nil
}

// UpdateLoginStreakLoginAtByID updates a user's login streak and login time
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

// UpdatePasswordByID updates a user's password and returns the updated user data
func (r *UserRepository) UpdatePasswordByID(ctx context.Context, id int, password string) (*dto.UserFullDTO, error) {
	u, err := r.client.User.
		UpdateOneID(id).
		SetPassword(password).
		Save(ctx)

	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.MapToUserFullDTOFromEnt(u), nil
}

// SetBlockUntilAndLevelAndReasonFromUser updates a user's block status information
func (r *UserRepository) SetBlockUntilAndLevelAndReasonFromUser(ctx context.Context, user *dto.UserDTO) error {
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

// Helper methods for error handling

// handleQueryError transforms Ent query errors into domain-specific errors
func (r *UserRepository) handleQueryError(err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return repositoryerrors.WrapErrUserNotFound(err)
	}

	return repositoryerrors.WrapUnexpectedError(err)
}

// handleUpdateError transforms Ent update errors into domain-specific errors
func (r *UserRepository) handleUpdateError(err error) error {
	if err == nil {
		return nil
	}

	return repositoryerrors.WrapUnexpectedError(err)
}

// handleConstraintError processes database constraint violation errors
func (r *UserRepository) handleConstraintError(err error) error {
	if err == nil {
		return nil
	}

	if !ent.IsConstraintError(err) {
		return repositoryerrors.WrapUnexpectedError(err)
	}

	switch {
	case strings.Contains(err.Error(), "username"):
		return repositoryerrors.WrapErrUserAlreadyExists(err)
	case strings.Contains(err.Error(), "hardware_id"):
		return repositoryerrors.WrapErrUserHwidConflict(err)
	default:
		return repositoryerrors.WrapUnexpectedError(err)
	}
}
