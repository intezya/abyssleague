package applicationservice

import (
	applicationerror "abysscore/internal/common/errors/application"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity"
	drivenports "abysscore/internal/domain/ports/driven"
	repositoryports "abysscore/internal/domain/repository"
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/internal/pkg/timeutils"
	"context"
	"github.com/intezya/pkglib/logger"
	"strings"
	"time"
)

// AuthenticationService handles user authentication operations
type AuthenticationService struct {
	userRepository       repositoryports.UserRepository
	mainWebsocketService drivenports.WebsocketService
	credentialsHelper    domainservice.CredentialsHelper
	tokenHelper          domainservice.TokenHelper
}

// NewAuthenticationService creates a new authentication service
func NewAuthenticationService(
	userRepository repositoryports.UserRepository,
	mainWebsocketService drivenports.WebsocketService,
	credentialsHelper domainservice.CredentialsHelper,
	tokenHelper domainservice.TokenHelper,
) *AuthenticationService {
	return &AuthenticationService{
		userRepository:       userRepository,
		mainWebsocketService: mainWebsocketService,
		credentialsHelper:    credentialsHelper,
		tokenHelper:          tokenHelper,
	}
}

// Register creates a new user account
func (a *AuthenticationService) Register(
	ctx context.Context,
	credentials *entity.CredentialsDTO,
) (
	*domainservice.AuthenticationResult,
	error,
) {
	// Encode sensitive data
	encodedPassword := tracer.TraceValue(
		ctx, "credentialsHelper.EncodePassword", func(ctx context.Context) string {
			return a.credentialsHelper.EncodePassword(credentials.Password)
		},
	)

	encodedHwid := tracer.TraceValue(
		ctx, "credentialsHelper.EncodeHardwareID", func(ctx context.Context) string {
			return a.credentialsHelper.EncodeHardwareID(credentials.Hwid)
		},
	)

	credentials.Password = encodedPassword
	credentials.Hwid = encodedHwid

	user, err := tracer.TraceFnWithResult(
		ctx, "userRepository.Create", func(ctx context.Context) (*entity.AuthenticationData, error) {
			return a.userRepository.Create(ctx, credentials)
		},
	)

	if err != nil {
		return nil, err // Username or hardware ID conflict
	}

	token := a.generateToken(ctx, user.TokenData())
	online := a.getOnlineCount(ctx)

	return domainservice.NewAuthenticationResult(token, nil, online), nil
}

// Authenticate validates user credentials and returns authentication result
func (a *AuthenticationService) Authenticate(ctx context.Context, credentials *entity.CredentialsDTO) (
	*domainservice.AuthenticationResult,
	error,
) {
	lowerUsername := strings.ToLower(credentials.Username)

	authentication, err := tracer.TraceFnWithResult(
		ctx,
		"userRepository.FindAuthenticationByLowerUsername",
		func(ctx context.Context) (*entity.AuthenticationData, error) {
			return a.userRepository.FindAuthenticationByLowerUsername(ctx, lowerUsername)
		},
	)

	if err != nil {
		return nil, err // User not found
	}

	// Verify password
	passwordOk := tracer.TraceValue(
		ctx, "authentication.ComparePassword", func(ctx context.Context) bool {
			return a.comparePasswords(authentication, credentials.Password)
		},
	)

	if !passwordOk {
		return nil, applicationerror.ErrWrongPassword
	}

	// Verify HWID
	hwidOk, needsUpdate := tracer.TraceValueValue(
		ctx, "authentication.CompareHWID", func(ctx context.Context) (bool, bool) {
			return authentication.CompareHWID(credentials.Hwid, a.credentialsHelper.VerifyHardwareID)
		},
	)

	if !hwidOk {
		return nil, applicationerror.ErrUserWrongHwid
	}

	// Update HWID if needed
	if needsUpdate {
		if err := a.updateHwid(ctx, authentication, credentials.Hwid); err != nil {
			return nil, err
		}
	}

	// Get full user data
	user, err := tracer.TraceFnWithResult(
		ctx, "userRepository.FindFullDTOById", func(ctx context.Context) (*dto.UserFullDTO, error) {
			return a.userRepository.FindFullDTOById(ctx, authentication.UserID())
		},
	)

	if err != nil {
		logger.Log.Warnw("Failed to retrieve full user data", "error", err, "userID", authentication.UserID())
		return nil, err
	}

	// Check if account is locked
	if a.isAccountLocked(user.UserDTO) {
		return nil, applicationerror.ErrAccountIsLocked
	}

	token := a.generateToken(ctx, authentication.TokenData())

	// Handle login streak and other post-login processing in background
	go a.processLoginStreakAndRewards(ctx, user.UserDTO)

	online := a.getOnlineCount(ctx)

	return domainservice.NewAuthenticationResult(token, user, online), nil
}

// ValidateToken validates the authentication token and returns user data
func (a *AuthenticationService) ValidateToken(ctx context.Context, token string) (*dto.UserDTO, error) {
	data, err := tracer.TraceFnWithResult(
		ctx, "tokenHelper.ValidateToken", func(ctx context.Context) (*entity.TokenData, error) {
			return a.tokenHelper.ValidateToken(token)
		},
	)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugw("authentication data received from token", "data", data)

	lowerUsername := strings.ToLower(data.Username)

	authentication, err := tracer.TraceFnWithResult(
		ctx,
		"userRepository.FindAuthenticationByLowerUsername",
		func(ctx context.Context) (*entity.AuthenticationData, error) {
			return a.userRepository.FindAuthenticationByLowerUsername(ctx, lowerUsername)
		},
	)

	if err != nil {
		return nil, err // User from token not found
	}

	// Verify HWID in token
	hwidOk, needsUpdate := tracer.TraceValueValue(
		ctx, "authentication.CompareHWID", func(ctx context.Context) (bool, bool) {
			return authentication.CompareHWID(data.Hwid, a.credentialsHelper.VerifyHardwareID)
		},
	)

	if !hwidOk || needsUpdate {
		return nil, applicationerror.ErrTokenHwidIsInvalid
	}

	user, err := tracer.TraceFnWithResult(
		ctx, "userRepository.FindDTOById", func(ctx context.Context) (*dto.UserDTO, error) {
			return a.userRepository.FindDTOById(ctx, authentication.UserID())
		},
	)

	if err != nil {
		logger.Log.Warnw(
			"Failed to retrieve user data during token validation",
			"error", err,
			"userID", authentication.UserID(),
		)
		return nil, err
	}

	if a.isAccountLocked(user) {
		return nil, applicationerror.ErrAccountIsLocked
	}

	return user, nil
}

// ChangePassword updates user password
func (a *AuthenticationService) ChangePassword(
	ctx context.Context,
	credentials *entity.ChangePasswordDTO,
) (*domainservice.AuthenticationResult, error) {
	lowerUsername := strings.ToLower(credentials.Username)

	authentication, err := tracer.TraceFnWithResult(
		ctx,
		"userRepository.FindAuthenticationByLowerUsername",
		func(ctx context.Context) (*entity.AuthenticationData, error) {
			return a.userRepository.FindAuthenticationByLowerUsername(ctx, lowerUsername)
		},
	)

	if err != nil {
		return nil, err // User not found
	}

	// Verify current password
	passwordOk := tracer.TraceValue(
		ctx, "authentication.ComparePassword", func(ctx context.Context) bool {
			return a.comparePasswords(authentication, credentials.OldPassword)
		},
	)

	if !passwordOk {
		return nil, applicationerror.ErrWrongPassword
	}

	// Encode new password
	encodedPassword := tracer.TraceValue(
		ctx, "credentialsHelper.EncodePassword", func(ctx context.Context) string {
			return a.credentialsHelper.EncodePassword(credentials.NewPassword)
		},
	)

	// Update password in database
	user, err := tracer.TraceFnWithResult(
		ctx,
		"userRepository.UpdatePasswordByID",
		func(ctx context.Context) (*dto.UserFullDTO, error) {
			return a.userRepository.UpdatePasswordByID(ctx, authentication.UserID(), encodedPassword)
		},
	)

	if err != nil {
		return nil, err
	}

	token := a.generateToken(ctx, authentication.TokenData())
	online := a.getOnlineCount(ctx)

	return domainservice.NewAuthenticationResult(token, user, online), nil
}

// Helper methods

// updateHwid updates the hardware ID for a user
func (a *AuthenticationService) updateHwid(ctx context.Context, auth *entity.AuthenticationData, rawHwid string) error {
	newHwid := tracer.TraceValue(
		ctx, "credentialsHelper.EncodeHardwareID", func(ctx context.Context) string {
			return a.credentialsHelper.EncodeHardwareID(rawHwid)
		},
	)

	auth.SetHWID(newHwid)

	err := tracer.TraceValue(
		ctx, "userRepository.UpdateHWIDByID", func(ctx context.Context) error {
			return a.userRepository.UpdateHWIDByID(ctx, auth.UserID(), newHwid)
		},
	)

	return err // Return hwid conflict if any
}

// processLoginStreakAndRewards handles post-login processing like login streaks and rewards
func (a *AuthenticationService) processLoginStreakAndRewards(ctx context.Context, user *dto.UserDTO) {
	// Only update streak if user hasn't logged in today
	if !timeutils.IsDayBeforeToday(user.LoginAt) {
		return
	}

	// Update login streak and login time
	user.LoginStreak++
	user.LoginAt = time.Now()

	// TODO: Implement logic to handle search block level decrement
	// TODO: Implement logic to add bonuses for user based on login streak

	err := tracer.TraceValue(
		ctx, "userRepository.SetLoginStreakLoginAtByID", func(ctx context.Context) error {
			return a.userRepository.SetLoginStreakLoginAtByID(ctx, user.ID, user.LoginStreak, user.LoginAt)
		},
	)

	if err != nil {
		logger.Log.Warnw("Failed to update login streak", "error", err, "userID", user.ID)
	}
}

// comparePasswords verifies the password
func (a *AuthenticationService) comparePasswords(authentication *entity.AuthenticationData, raw string) bool {
	return authentication.ComparePassword(raw, a.credentialsHelper.VerifyPassword)
}

// isAccountLocked checks if user account is locked
func (a *AuthenticationService) isAccountLocked(user *dto.UserDTO) bool {
	return user.AccountBlockedUntil.After(time.Now())
}

// generateToken creates an authentication token
func (a *AuthenticationService) generateToken(ctx context.Context, tokenData *entity.TokenData) string {
	return tracer.TraceValue(
		ctx, "tokenHelper.TokenGenerator", func(ctx context.Context) string {
			return a.tokenHelper.TokenGenerator(tokenData)
		},
	)
}

// getOnlineCount retrieves the number of online users
func (a *AuthenticationService) getOnlineCount(ctx context.Context) int {
	return tracer.TraceValue(
		ctx, "mainWebsocketService.GetOnlineSoft", func(ctx context.Context) int {
			res := a.mainWebsocketService.GetOnlineSoft(ctx)
			return int(res.Online)
		},
	)
}
