package applicationservice

import (
	"context"
	"strings"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/clients"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/userentity"
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
	"github.com/intezya/abyssleague/services/abysscore/pkg/timeutils"
	"github.com/intezya/pkglib/logger"
)

// AuthenticationService handles user authentication operations.
type AuthenticationService struct {
	authRepo             repositoryports.AuthenticationRepository
	userRepo             repositoryports.UserRepository
	websocketClient      clients.WebsocketMessagingClient
	credentialsHelper    domainservice.CredentialsHelper
	tokenHelper          domainservice.TokenHelper
	bannedHardwareIDRepo repositoryports.BannedHardwareIDRepository
}

// NewAuthenticationService creates a new authentication service with dependency injection.
func NewAuthenticationService(
	authRepo repositoryports.AuthenticationRepository,
	userRepo repositoryports.UserRepository,
	websocketClient clients.WebsocketMessagingClient,
	credentialsHelper domainservice.CredentialsHelper,
	tokenHelper domainservice.TokenHelper,
	bannedHardwareIDRepo repositoryports.BannedHardwareIDRepository,
) *AuthenticationService {
	return &AuthenticationService{
		authRepo:             authRepo,
		userRepo:             userRepo,
		websocketClient:      websocketClient,
		credentialsHelper:    credentialsHelper,
		tokenHelper:          tokenHelper,
		bannedHardwareIDRepo: bannedHardwareIDRepo,
	}
}

// Register creates a new user account.
func (s *AuthenticationService) Register(
	ctx context.Context,
	credentials *dto.CredentialsDTO,
) (*domainservice.AuthenticationResult, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.Register")
	defer span.End()

	s.prepareCredentials(ctx, credentials)

	tx, err := s.authRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}

	result, err := persistence.WithTxResultTx(ctx, tx, func(tx *ent.Tx) (*entity.AuthenticationData, error) {
		// Check if hardware ID is banned
		hardwareIDBanned, err := s.bannedHardwareIDRepo.TxFindByHardwareID(ctx, tx, credentials.HardwareID)
		if hardwareIDBanned != nil {
			return nil, apperrors.ErrHardwareIDBanned(hardwareIDBanned.BanReason)
		}

		// Create user
		user, err := s.userRepo.TxCreate(ctx, tx, credentials)
		if err != nil {
			return nil, err
		}

		return user, nil
	})
	if err != nil {
		return nil, err // Username or hardware ID conflict
	}

	return s.createAuthResult(ctx, result.TokenData(), nil), nil
}

// Authenticate validates user credentials and returns authentication result.
func (s *AuthenticationService) Authenticate(
	ctx context.Context,
	credentials *dto.CredentialsDTO,
) (*domainservice.AuthenticationResult, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.Authenticate")
	defer span.End()

	// Encode hardware ID for comparison
	credentials.HardwareID = s.encodeHardwareID(ctx, credentials.HardwareID)

	tx, err := s.authRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}

	authentication, user, err := persistence.WithTxResult2Tx(
		ctx,
		tx,
		func(tx *ent.Tx) (*entity.AuthenticationData, *dto.UserFullDTO, error) {
			// Find authentication data by username
			authentication, err := s.findAuthByUsername(ctx, tx, credentials.Username)
			if err != nil {
				return nil, nil, err
			}

			// Verify password
			if !s.verifyPassword(ctx, authentication, credentials.Password) {
				return nil, nil, apperrors.ErrWrongPassword
			}

			// Verify and update hardware ID if necessary
			if err := s.verifyAndUpdateHardwareID(ctx, tx, authentication, credentials.HardwareID); err != nil {
				return nil, nil, err
			}

			// Check if account is locked
			if authentication.IsAccountLocked() {
				return nil, nil, apperrors.ErrAccountIsLocked(authentication.BlockReason())
			}

			// Get full user data
			user, err := s.userRepo.FindFullDTOById(ctx, authentication.UserID())
			if err != nil {
				logger.Log.Warnw(
					"Failed to retrieve full user data",
					"error", err,
					"userID", authentication.UserID(),
				)
				return nil, nil, err
			}

			// Process post-login tasks
			s.processPostLoginTasks(ctx, tx, user.UserDTO)

			return authentication, user, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return s.createAuthResult(ctx, authentication.TokenData(), user), nil
}

// ValidateToken validates the authentication token and returns user data.
func (s *AuthenticationService) ValidateToken(
	ctx context.Context,
	token string,
) (*dto.UserDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.ValidateToken")
	defer span.End()

	// Validate token and extract data
	tokenData, err := s.tokenHelper.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugw("authentication data received from token", "data", tokenData)

	tx, err := s.authRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}

	user, err := persistence.WithTxResultTx(ctx, tx, func(tx *ent.Tx) (*dto.UserDTO, error) {
		// Find authentication data by username from token
		authentication, err := s.findAuthByUsername(ctx, tx, tokenData.Username)
		if err != nil {
			return nil, err
		}

		// Verify hardware ID from token
		hardwareIDOk, needsUpdate := authentication.CompareHardwareID(
			tokenData.Hwid,
			s.credentialsHelper.VerifyTokenHardwareID,
		)

		if !hardwareIDOk || needsUpdate {
			return nil, apperrors.ErrTokenHardwareIDIsInvalid
		}

		// Check if account is locked
		if authentication.IsAccountLocked() {
			return nil, apperrors.ErrAccountIsLocked(authentication.BlockReason())
		}

		// Get user data
		user, err := s.userRepo.FindDTOById(ctx, authentication.UserID())
		if err != nil {
			logger.Log.Warnw(
				"Failed to retrieve user data during token validation",
				"error", err,
				"userID", authentication.UserID(),
			)
			return nil, err
		}

		return user, nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword updates user password.
func (s *AuthenticationService) ChangePassword(
	ctx context.Context,
	credentials *dto.ChangePasswordDTO,
) (*domainservice.AuthenticationResult, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.ChangePassword")
	defer span.End()

	tx, err := s.authRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}

	authentication, user, err := persistence.WithTxResult2Tx(
		ctx,
		tx,
		func(tx *ent.Tx) (*entity.AuthenticationData, *dto.UserFullDTO, error) {
			// Find authentication by username
			authentication, err := s.findAuthByUsername(ctx, tx, credentials.Username)
			if err != nil {
				return nil, nil, err
			}

			// Verify old password
			if !s.verifyPassword(ctx, authentication, credentials.OldPassword) {
				return nil, nil, apperrors.ErrWrongPassword
			}

			// Encode new password
			encodedPassword := s.encodePassword(ctx, credentials.NewPassword)

			// Update password and get user data
			user, err := s.authRepo.TxUpdatePasswordByID(ctx, tx, authentication.UserID(), encodedPassword)
			if err != nil {
				return nil, nil, err
			}

			return authentication, user, nil
		})

	if err != nil {
		return nil, err
	}

	return s.createAuthResult(ctx, authentication.TokenData(), user), nil
}

// Helper methods

// prepareCredentials encodes password and HWID.
func (s *AuthenticationService) prepareCredentials(
	ctx context.Context,
	credentials *dto.CredentialsDTO,
) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.prepareCredentials")
	defer span.End()

	credentials.Password = s.encodePassword(ctx, credentials.Password)
	credentials.HardwareID = s.encodeHardwareID(ctx, credentials.HardwareID)
}

// findAuthByUsername finds authentication data by username.
func (s *AuthenticationService) findAuthByUsername(
	ctx context.Context,
	tx *ent.Tx,
	username string,
) (*entity.AuthenticationData, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.findAuthByUsername")
	defer span.End()

	lowerUsername := strings.ToLower(username)
	return s.authRepo.TxFindAuthenticationByLowerUsername(ctx, tx, lowerUsername)
}

// verifyPassword checks if the provided password matches the stored one.
func (s *AuthenticationService) verifyPassword(
	ctx context.Context,
	auth *entity.AuthenticationData,
	password string,
) bool {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.verifyPassword")
	defer span.End()

	return auth.ComparePassword(password, s.credentialsHelper.VerifyPassword)
}

// verifyAndUpdateHardwareID validates HardwareID and updates it if necessary.
func (s *AuthenticationService) verifyAndUpdateHardwareID(
	ctx context.Context,
	tx *ent.Tx,
	auth *entity.AuthenticationData,
	hardwareID string,
) error {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.verifyAndUpdateHardwareID")
	defer span.End()

	hardwareIDOk, needsUpdate := auth.CompareHardwareID(hardwareID, s.credentialsHelper.VerifyHardwareID)

	if !hardwareIDOk {
		return apperrors.ErrUserWrongHardwareID
	}

	if needsUpdate {
		auth.SetHardwareID(hardwareID)

		err := s.authRepo.TxUpdateHardwareIDByID(ctx, tx, auth.UserID(), hardwareID)
		if err != nil {
			return err
		}
	}

	return nil
}

// encodePassword encodes a raw password.
func (s *AuthenticationService) encodePassword(ctx context.Context, rawPassword string) string {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.encodePassword")
	defer span.End()

	return s.credentialsHelper.EncodePassword(rawPassword)
}

// encodeHardwareID encodes a raw hardware ID.
func (s *AuthenticationService) encodeHardwareID(ctx context.Context, rawHwid string) string {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.encodeHardwareID")
	defer span.End()

	return s.credentialsHelper.EncodeHardwareID(rawHwid)
}

// createAuthResult creates authentication result with token and online count.
func (s *AuthenticationService) createAuthResult(
	ctx context.Context,
	tokenData *entity.TokenData,
	user *dto.UserFullDTO,
) *domainservice.AuthenticationResult {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.createAuthResult")
	defer span.End()

	token := s.generateToken(ctx, tokenData)
	online := s.getOnlineCount(ctx)

	return domainservice.NewAuthenticationResult(token, user, online)
}

// processPostLoginTasks handles all post-login actions.
func (s *AuthenticationService) processPostLoginTasks(
	ctx context.Context,
	tx *ent.Tx,
	user *dto.UserDTO,
) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.processPostLoginTasks")
	defer span.End()

	s.processLoginStreakAndRewards(ctx, tx, user)
	s.processBanDecrementAfterLogin(ctx, tx, user)
}

// processLoginStreakAndRewards handles post-login processing like login streaks and rewards.
func (s *AuthenticationService) processLoginStreakAndRewards(
	ctx context.Context,
	tx *ent.Tx,
	user *dto.UserDTO,
) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.processLoginStreakAndRewards")
	defer span.End()

	// Only update streak if user hasn't logged in today
	if !timeutils.IsDayBeforeToday(user.LoginAt) {
		return
	}

	user.LoginStreak++
	user.LoginAt = time.Now()

	// TODO: Implement logic to handle search block level decrement
	// TODO: Implement logic to add bonuses for user based on login streak

	err := s.authRepo.TxUpdateLoginStreakLoginAtByID(
		ctx,
		tx,
		user.ID,
		user.LoginStreak,
		user.LoginAt,
	)

	if err != nil {
		logger.Log.Warnw("Failed to update login streak", "error", err, "userID", user.ID)
	}
}

// processBanDecrementAfterLogin decrements ban levels if needed.
func (s *AuthenticationService) processBanDecrementAfterLogin(
	ctx context.Context,
	tx *ent.Tx,
	user *dto.UserDTO,
) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.processBanDecrementAfterLogin")
	defer span.End()

	var now = time.Now()

	// Process account block decrement
	if user.AccountBlockedUntil != nil &&
		user.AccountBlockedUntil.Add(userentity.AccountBlockDecrementTime).Before(now) {
		if user.AccountBlockedLevel > 0 {
			user.AccountBlockedLevel--
		}
		user.AccountBlockedUntil = nil
		user.AccountBlockReason = nil
	}

	// Process search block decrement
	if user.SearchBlockedUntil != nil &&
		user.SearchBlockedUntil.Add(userentity.SearchBlockDecrementTime).Before(now) {
		if user.SearchBlockedLevel > 0 {
			user.SearchBlockedLevel--
		}
		user.SearchBlockedUntil = nil
		user.SearchBlockReason = nil
	}

	err := s.authRepo.TxSetBlockUntilAndLevelAndReasonFromUser(ctx, tx, user)
	if err != nil {
		logger.Log.Errorw(
			"Failed to update block settings for user",
			"error", err,
			"userID", user.ID,
		)
	}
}

// generateToken creates an authentication token.
func (s *AuthenticationService) generateToken(
	ctx context.Context,
	tokenData *entity.TokenData,
) string {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.generateToken")
	defer span.End()

	return s.tokenHelper.TokenGenerator(tokenData)
}

// getOnlineCount retrieves the number of online users.
func (s *AuthenticationService) getOnlineCount(ctx context.Context) int {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.getOnlineCount")
	defer span.End()

	res, err := s.websocketClient.GetOnline(ctx)
	if err != nil {
		logger.Log.Debugw("Failed to get online count", "error", err)
		return 0
	}

	return res
}
