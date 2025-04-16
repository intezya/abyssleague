package domainservice

import (
	"abysscore/internal/domain/entity/userentity"
	"context"
)

type AuthenticationResult struct {
	Token string                  `json:"token,omitempty"`
	User  *userentity.UserFullDTO `json:"user,omitempty"`
}

func NewAuthenticationResult(token string, user *userentity.UserFullDTO) *AuthenticationResult {
	return &AuthenticationResult{Token: token, User: user}
}

type AuthenticationService interface {
	Register(ctx context.Context, credentials *userentity.CredentialsDTO) (*AuthenticationResult, error)
	Authenticate(ctx context.Context, credentials *userentity.CredentialsDTO) (*AuthenticationResult, error)
	ValidateToken(ctx context.Context, token string) (*userentity.UserDTO, error)
}

type TokenHelper interface {
	TokenGenerator(tokenData *userentity.TokenData) string
	ValidateToken(token string) (*userentity.TokenData, error)
}

type CredentialsHelper interface {
	EncodePassword(raw string) (hash string)
	VerifyPassword(raw, hash string) bool
	EncodeHardwareID(raw string) (hash string)
	VerifyHardwareID(raw, hash string) bool
}
