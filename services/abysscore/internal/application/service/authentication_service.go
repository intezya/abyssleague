package applicationservice

import (
	applicationerror "abysscore/internal/common/errors/application"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity"
	"abysscore/internal/domain/entity/userentity"
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
	credentials *dto.CredentialsDTO,
) (*domainservice.AuthenticationResult, error) {
	a.prepareCredentials(ctx, credentials)

	user, err := tracer.TraceFnWithResult(
		ctx, "userRepository.Create", func(ctx context.Context) (*entity.AuthenticationData, error) {
			return a.userRepository.Create(ctx, credentials)
		},
	)
	if err != nil {
		return nil, err // Username or hardware ID conflict
	}

	return a.createAuthResult(ctx, user.TokenData(), nil), nil
}

// Authenticate validates user credentials and returns authentication result
func (a *AuthenticationService) Authenticate(ctx context.Context, credentials *dto.CredentialsDTO) (
	*domainservice.AuthenticationResult,
	error,
) {
	authentication, err := a.findAuthByUsername(ctx, credentials.Username)
	if err != nil {
		return nil, err // User not found
	}

	if !a.verifyPassword(ctx, authentication, credentials.Password) {
		return nil, applicationerror.ErrWrongPassword
	}

	if err := a.verifyAndUpdateHWID(ctx, authentication, credentials.Hwid); err != nil {
		return nil, err
	}

	if authentication.IsAccountLocked() {
		return nil, applicationerror.ErrAccountIsLocked(authentication.BlockReason())
	}

	user, err := tracer.TraceFnWithResult(
		ctx, "userRepository.FindFullDTOById", func(ctx context.Context) (*dto.UserFullDTO, error) {
			return a.userRepository.FindFullDTOById(ctx, authentication.UserID())
		},
	)
	if err != nil {
		logger.Log.Warnw("Failed to retrieve full user data", "error", err, "userID", authentication.UserID())
		return nil, err
	}

	// Async login processing (statistics, bonuses, etc)
	go a.processPostLoginTasks(ctx, user.UserDTO)

	return a.createAuthResult(ctx, authentication.TokenData(), user), nil
}

// ValidateToken validates the authentication token and returns user data
func (a *AuthenticationService) ValidateToken(ctx context.Context, token string) (*dto.UserDTO, error) {
	// Валидация токена
	tokenData, err := tracer.TraceFnWithResult(
		ctx, "tokenHelper.ValidateToken", func(ctx context.Context) (*entity.TokenData, error) {
			return a.tokenHelper.ValidateToken(token)
		},
	)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugw("authentication data received from token", "data", tokenData)

	authentication, err := a.findAuthByUsername(ctx, tokenData.Username)
	if err != nil {
		return nil, err
	}

	hwidOk, needsUpdate := tracer.Trace2(
		ctx, "authentication.CompareHWID", func(ctx context.Context) (bool, bool) {
			return authentication.CompareHWID(tokenData.Hwid, a.credentialsHelper.VerifyHardwareID)
		},
	)
	if !hwidOk || needsUpdate {
		return nil, applicationerror.ErrTokenHwidIsInvalid
	}

	if authentication.IsAccountLocked() {
		return nil, applicationerror.ErrAccountIsLocked(authentication.BlockReason())
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

	return user, nil
}

// ChangePassword updates user password
func (a *AuthenticationService) ChangePassword(
	ctx context.Context,
	credentials *dto.ChangePasswordDTO,
) (*domainservice.AuthenticationResult, error) {
	authentication, err := a.findAuthByUsername(ctx, credentials.Username)
	if err != nil {
		return nil, err // User not found
	}

	if !a.verifyPassword(ctx, authentication, credentials.OldPassword) {
		return nil, applicationerror.ErrWrongPassword
	}

	encodedPassword := a.encodePassword(ctx, credentials.NewPassword)

	user, err := tracer.TraceFnWithResult(
		ctx, "userRepository.UpdatePasswordByID", func(ctx context.Context) (*dto.UserFullDTO, error) {
			return a.userRepository.UpdatePasswordByID(ctx, authentication.UserID(), encodedPassword)
		},
	)
	if err != nil {
		return nil, err
	}

	return a.createAuthResult(ctx, authentication.TokenData(), user), nil
}

/*
	Helper methods
*/

// prepareCredentials encodes password and HWID
func (a *AuthenticationService) prepareCredentials(ctx context.Context, credentials *dto.CredentialsDTO) {
	credentials.Password = a.encodePassword(ctx, credentials.Password)
	credentials.Hwid = a.encodeHWID(ctx, credentials.Hwid)
}

// findAuthByUsername finds authentication data by username
func (a *AuthenticationService) findAuthByUsername(ctx context.Context, username string) (
	*entity.AuthenticationData,
	error,
) {
	lowerUsername := strings.ToLower(username)
	return tracer.TraceFnWithResult(
		ctx, "userRepository.FindAuthenticationByLowerUsername",
		func(ctx context.Context) (*entity.AuthenticationData, error) {
			return a.userRepository.FindAuthenticationByLowerUsername(ctx, lowerUsername)
		},
	)
}

// verifyPassword checks if the provided password matches the stored one
func (a *AuthenticationService) verifyPassword(
	ctx context.Context,
	auth *entity.AuthenticationData,
	password string,
) bool {
	return tracer.Trace1(
		ctx, "authentication.ComparePassword", func(ctx context.Context) bool {
			return auth.ComparePassword(password, a.credentialsHelper.VerifyPassword)
		},
	)
}

// verifyAndUpdateHWID validates HWID and updates it if necessary
func (a *AuthenticationService) verifyAndUpdateHWID(
	ctx context.Context,
	auth *entity.AuthenticationData,
	hwid string,
) error {
	// Проверка HWID
	hwidOk, needsUpdate := tracer.Trace2(
		ctx, "authentication.CompareHWID", func(ctx context.Context) (bool, bool) {
			return auth.CompareHWID(hwid, a.credentialsHelper.VerifyHardwareID)
		},
	)

	if !hwidOk {
		return applicationerror.ErrUserWrongHwid
	}

	// Обновление HWID если необходимо
	if needsUpdate {
		if err := a.updateHwid(ctx, auth, hwid); err != nil {
			return err
		}
	}

	return nil
}

// updateHwid updates the hardware ID for a user
func (a *AuthenticationService) updateHwid(ctx context.Context, auth *entity.AuthenticationData, rawHwid string) error {
	newHwid := a.encodeHWID(ctx, rawHwid)
	auth.SetHWID(newHwid)

	return tracer.Trace1(
		ctx, "userRepository.UpdateHWIDByID", func(ctx context.Context) error {
			return a.userRepository.UpdateHWIDByID(ctx, auth.UserID(), newHwid)
		},
	)
}

// encodePassword encodes a raw password
func (a *AuthenticationService) encodePassword(ctx context.Context, rawPassword string) string {
	return tracer.Trace1(
		ctx, "credentialsHelper.EncodePassword", func(ctx context.Context) string {
			return a.credentialsHelper.EncodePassword(rawPassword)
		},
	)
}

// encodeHWID encodes a raw hardware ID
func (a *AuthenticationService) encodeHWID(ctx context.Context, rawHwid string) string {
	return tracer.Trace1(
		ctx, "credentialsHelper.EncodeHardwareID", func(ctx context.Context) string {
			return a.credentialsHelper.EncodeHardwareID(rawHwid)
		},
	)
}

// createAuthResult creates authentication result with token and online count
func (a *AuthenticationService) createAuthResult(
	ctx context.Context,
	tokenData *entity.TokenData,
	user *dto.UserFullDTO,
) *domainservice.AuthenticationResult {
	token := a.generateToken(ctx, tokenData)
	online := a.getOnlineCount(ctx)
	return domainservice.NewAuthenticationResult(token, user, online)
}

// processPostLoginTasks handles all post-login actions
func (a *AuthenticationService) processPostLoginTasks(ctx context.Context, user *dto.UserDTO) {
	a.processLoginStreakAndRewards(ctx, user)
	a.processBanDecrementAfterLogin(ctx, user)
}

// processLoginStreakAndRewards handles post-login processing like login streaks and rewards
func (a *AuthenticationService) processLoginStreakAndRewards(ctx context.Context, user *dto.UserDTO) {
	// Only update streak if user hasn't logged in today
	if !timeutils.IsDayBeforeToday(user.LoginAt) {
		return
	}

	user.LoginStreak++
	user.LoginAt = time.Now()

	// TODO: Implement logic to handle search block level decrement
	// TODO: Implement logic to add bonuses for user based on login streak

	err := tracer.Trace1(
		ctx, "userRepository.UpdateLoginStreakLoginAtByID", func(ctx context.Context) error {
			return a.userRepository.UpdateLoginStreakLoginAtByID(ctx, user.ID, user.LoginStreak, user.LoginAt)
		},
	)

	if err != nil {
		logger.Log.Warnw("Failed to update login streak", "error", err, "userID", user.ID)
	}
}

// processBanDecrementAfterLogin decrements ban levels if needed
func (a *AuthenticationService) processBanDecrementAfterLogin(ctx context.Context, user *dto.UserDTO) {
	if user.AccountBlockedUntil != nil && user.AccountBlockedUntil.Add(userentity.AccountBlockDecrementTime).Before(time.Now()) {
		if user.AccountBlockedLevel > 0 {
			user.AccountBlockedLevel--
		}
		user.AccountBlockedUntil = nil
		user.AccountBlockReason = nil
	}

	if user.SearchBlockedUntil != nil && user.SearchBlockedUntil.Add(userentity.SearchBlockDecrementTime).Before(time.Now()) {
		if user.SearchBlockedLevel > 0 {
			user.SearchBlockedLevel--
		}
		user.SearchBlockedUntil = nil
		user.SearchBlockReason = nil
	}

	err := tracer.TraceFn(ctx, "userRepository.SetBlockUntilAndLevelAndReasonFromUser", func(ctx context.Context) error {
		return a.userRepository.SetBlockUntilAndLevelAndReasonFromUser(ctx, user)
	})

	if err != nil {
		logger.Log.Errorw("Failed to update block until, level, user", "error", err, "userID", user.ID)
	}
}

// generateToken creates an authentication token
func (a *AuthenticationService) generateToken(ctx context.Context, tokenData *entity.TokenData) string {
	return tracer.Trace1(
		ctx, "tokenHelper.TokenGenerator", func(ctx context.Context) string {
			return a.tokenHelper.TokenGenerator(tokenData)
		},
	)
}

// getOnlineCount retrieves the number of online users
func (a *AuthenticationService) getOnlineCount(ctx context.Context) int {
	return tracer.Trace1(
		ctx, "mainWebsocketService.GetOnlineSoft", func(ctx context.Context) int {
			res, err := a.mainWebsocketService.GetOnline(ctx)

			if err != nil {
				return 0
			}

			return int(res.Online)
		},
	)
}
