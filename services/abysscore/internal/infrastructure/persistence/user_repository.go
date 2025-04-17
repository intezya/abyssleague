package persistence

import (
	repositoryerrors "abysscore/internal/common/errors/repository"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/userentity"
	"abysscore/internal/infrastructure/ent"
	"abysscore/internal/infrastructure/ent/user"
	"context"
	"strings"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (r *UserRepository) Create(credentials *userentity.CredentialsDTO) (*userentity.AuthenticationData, error) {
	u, err := r.client.User.
		Create().
		SetUsername(credentials.Username).
		SetLowerUsername(strings.ToLower(credentials.Username)).
		SetPassword(credentials.Password).
		SetHardwareID(credentials.Hwid).
		Save(context.Background())

	if err != nil && ent.IsConstraintError(err) {
		if strings.Contains(err.Error(), "username") {
			return nil, repositoryerrors.WrapErrUserAlreadyExists(err)
		} else if strings.Contains(err.Error(), "hardware_id") {
			return nil, repositoryerrors.WrapErrUserHwidConflict(err)
		}
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	return dto.MapToAuthenticationDataFromEnt(u), nil
}

func (r *UserRepository) FindDTOById(id int) (*dto.UserDTO, error) {
	ctx := context.Background()

	u, err := r.client.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, repositoryerrors.WrapErrUserNotFound(err)
		}
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	return dto.MapToUserDTOFromEnt(u), nil
}

func (r *UserRepository) FindFullDTOById(id int) (*dto.UserFullDTO, error) {
	ctx := context.Background()

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
		if ent.IsNotFound(err) {
			return nil, repositoryerrors.WrapErrUserNotFound(err)
		}
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	return dto.MapToFullDTOFromEnt(u), nil
}

func (r *UserRepository) FindAuthenticationByLowerUsername(lowerUsername string) (
	*userentity.AuthenticationData,
	error,
) {
	ctx := context.Background()

	u, err := r.client.User.
		Query().
		Where(user.LowerUsernameEQ(lowerUsername)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, repositoryerrors.WrapErrUserNotFound(err)
		}
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	return dto.MapToAuthenticationDataFromEnt(u), nil
}

func (r *UserRepository) UpdateHWIDByID(id int, hwid string) error {
	ctx := context.Background()

	_, err := r.client.User.
		UpdateOneID(id).
		SetHardwareID(hwid).
		Save(ctx)

	if err != nil && ent.IsConstraintError(err) {
		if strings.Contains(err.Error(), "hwid") {
			return repositoryerrors.WrapErrUserHwidConflict(err)
		}
		return repositoryerrors.WrapUnexpectedError(err)
	}

	return nil
}
