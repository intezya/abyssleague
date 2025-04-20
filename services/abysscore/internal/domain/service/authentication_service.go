package domainservice

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity"
	"context"
)

type AuthenticationResult struct {
	Token       string           `json:"token,omitempty"`
	User        *dto.UserFullDTO `json:"user,omitempty"`
	OnlineCount int              `json:"online_count"`
}

func NewAuthenticationResult(token string, user *dto.UserFullDTO, onlineCount int) *AuthenticationResult {
	return &AuthenticationResult{Token: token, User: user, OnlineCount: onlineCount}
}

type AuthenticationService interface {
	Register(ctx context.Context, credentials *dto.CredentialsDTO) (*AuthenticationResult, error)
	Authenticate(ctx context.Context, credentials *dto.CredentialsDTO) (*AuthenticationResult, error)
	ValidateToken(ctx context.Context, token string) (*dto.UserDTO, error)
	ChangePassword(ctx context.Context, credentials *dto.ChangePasswordDTO) (*AuthenticationResult, error)
}

type TokenHelper interface {
	TokenGenerator(tokenData *entity.TokenData) string
	ValidateToken(token string) (*entity.TokenData, error)
}

type CredentialsHelper interface {
	EncodePassword(raw string) (hash string)
	VerifyPassword(raw, hash string) bool
	EncodeHardwareID(raw string) (hash string)
	VerifyHardwareID(raw, hash string) bool
}
