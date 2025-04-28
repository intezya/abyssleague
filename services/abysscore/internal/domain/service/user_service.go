package domainservice

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
)

type UserService interface {
	SendCodeForEmailLink(ctx context.Context, user *dto.UserDTO, email string) error

	EnterCodeForEmailLink(
		ctx context.Context,
		user *dto.UserDTO,
		verificationCode string,
	) (*dto.UserDTO, error)
}
