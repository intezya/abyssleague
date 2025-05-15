package domainservice

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
)

type AuthenticationEventService interface {
	HandleRegistration(ctx context.Context, user *dto.UserDTO)
	HandleLogin(ctx context.Context, user *dto.UserDTO)
}
