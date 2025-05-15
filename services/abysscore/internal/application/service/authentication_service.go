package applicationservice

import (
	"context"
	"strings"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/clients"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
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
	eventService         domainservice.AuthenticationEventService
}

// NewAuthenticationService creates a new authentication service with dependency injection.
func NewAuthenticationService(
	authRepo repositoryports.AuthenticationRepository,
	userRepo repositoryports.UserRepository,
	websocketClient clients.WebsocketMessagingClient,
	credentialsHelper domainservice.CredentialsHelper,
	tokenHelper domainservice.TokenHelper,
	bannedHardwareIDRepo repositoryports.BannedHardwareIDRepository,
	eventService domainservice.AuthenticationEventService,
) *AuthenticationService {
	return &AuthenticationService{
		authRepo:             authRepo,
		userRepo:             userRepo,
		websocketClient:      websocketClient,
		credentialsHelper:    credentialsHelper,
		tokenHelper:          tokenHelper,
		bannedHardwareIDRepo: bannedHardwareIDRepo,
		eventService:         eventService,
	}
}

// Register creates a new user account.
func (s *AuthenticationService) Register(
	ctx context.Context,
	credentials *dto.CredentialsDTO,
) (*domainservice.AuthenticationResult, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.Register")
	defer span.End()

	tx, err := s.authRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}

	s.encryptCredentials(ctx, credentials)

	user, err := persistence.WithTxResultTx(
		ctx, tx, func(tx *ent.Tx) (*dto.UserDTO, error) {
			if err := s.checkHardwareIDBanned(ctx, tx, credentials.HardwareID); err != nil {
				return nil, err
			}

			user, err := s.userRepo.TxCreate(ctx, tx, credentials)
			if err != nil {
				return nil, err
			}

			return user, nil
		},
	)
	if err != nil {
		return nil, err
	}

	s.eventService.HandleRegistration(ctx, user)

	return s.createAuthResult(ctx, &dto.UserFullDTO{UserDTO: user}), nil
}

// Authenticate validates user credentials and returns authentication result.
func (s *AuthenticationService) Authenticate(
	ctx context.Context,
	credentials *dto.CredentialsDTO,
) (*domainservice.AuthenticationResult, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.Authenticate")
	defer span.End()

	tx, err := s.authRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}

	user, err := persistence.WithTxResultTx(
		ctx, tx, func(tx *ent.Tx) (*dto.UserFullDTO, error) {
			user, err := s.findUserFullDTOByUsername(ctx, tx, credentials.Username)
			if err != nil {
				return nil, err
			}

			if err := s.verifyAndUpdateHardwareID(ctx, tx, user, credentials.HardwareID); err != nil {
				return nil, err
			}

			if !s.verifyPassword(ctx, user.Password, credentials.Password) {
				return nil, apperrors.ErrWrongPassword
			}

			if s.isAccountLocked(user.UserDTO) {
				return nil, apperrors.ErrAccountIsLocked(user.AccountBlockReason)
			}

			return user, nil
		},
	)
	if err != nil {
		return nil, err
	}

	s.eventService.HandleLogin(ctx, user.UserDTO)

	return s.createAuthResult(ctx, user), nil
}

// ValidateToken validates the authentication token and returns user data.
func (s *AuthenticationService) ValidateToken(
	ctx context.Context,
	token string,
) (*dto.UserDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.ValidateToken")
	defer span.End()

	tokenData, err := s.tokenHelper.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugw("authentication data received from token", "data", tokenData)

	tx, err := s.authRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}

	user, err := persistence.WithTxResultTx(
		ctx, tx, func(tx *ent.Tx) (*dto.UserDTO, error) {
			// Find user auth by username from token
			user, err := s.findUserDTOByUsername(ctx, tx, tokenData.Username)
			if err != nil {
				return nil, err
			}

			// Verify hardware ID from token
			err = s.verifyTokenHardwareID(ctx, tx, user.HardwareID, tokenData.HardwareID)
			if err != nil {
				return nil, err
			}

			// Check if account is locked
			if s.isAccountLocked(user) {
				return nil, apperrors.ErrAccountIsLocked(user.AccountBlockReason)
			}

			return user, nil
		},
	)
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

	userAuth, err := persistence.WithTxResultTx(
		ctx,
		tx,
		func(tx *ent.Tx) (*dto.UserFullDTO, error) {
			// Find user auth by username
			userAuth, err := s.findUserFullDTOByUsername(ctx, tx, credentials.Username)
			if err != nil {
				return nil, err
			}

			// Verify old password
			if !s.verifyPassword(ctx, userAuth.Password, credentials.OldPassword) {
				return nil, apperrors.ErrWrongPassword
			}

			// Encode new password
			encodedPassword := s.encodePassword(ctx, credentials.NewPassword)

			// Update password
			err = s.authRepo.TxUpdatePasswordByID(ctx, tx, userAuth.ID, encodedPassword)
			if err != nil {
				return nil, err
			}

			// Update user auth object with new password
			userAuth.Password = encodedPassword

			return userAuth, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return s.createAuthResult(ctx, userAuth), nil
}

// Helper methods

// findUserFullDTOByUsername finds user with authentication data by username.
func (s *AuthenticationService) findUserFullDTOByUsername(
	ctx context.Context,
	tx *ent.Tx,
	username string,
) (*dto.UserFullDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.findUserFullDTOByUsername")
	defer span.End()

	lowerUsername := strings.ToLower(username)

	return s.userRepo.TxFindFullDTOByLowerUsername(ctx, tx, lowerUsername)
}

// findUserDTOByUsername finds user with authentication data by username.
func (s *AuthenticationService) findUserDTOByUsername(
	ctx context.Context,
	tx *ent.Tx,
	username string,
) (*dto.UserDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.findUserDTOByUsername")
	defer span.End()

	lowerUsername := strings.ToLower(username)

	return s.userRepo.TxFindDTOByLowerUsername(ctx, tx, lowerUsername)
}

// encryptCredentials encodes password and HWID.
func (s *AuthenticationService) encryptCredentials(
	ctx context.Context,
	credentials *dto.CredentialsDTO,
) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.encryptCredentials")
	defer span.End()

	credentials.Password = s.encodePassword(ctx, credentials.Password)
	credentials.HardwareID = s.encodeHardwareID(ctx, credentials.HardwareID)
}

// verifyPassword checks if the provided password matches the stored one.
func (s *AuthenticationService) verifyPassword(
	ctx context.Context,
	hashedPassword string,
	rawPassword string,
) bool {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.verifyPassword")
	defer span.End()

	return s.credentialsHelper.VerifyPassword(rawPassword, hashedPassword)
}

// verifyTokenHardwareID checks if the token's HWID matches the stored one.
func (s *AuthenticationService) verifyTokenHardwareID(
	ctx context.Context,
	tx *ent.Tx,
	storedHardwareID *string,
	tokenHardwareID string,
) error {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.verifyTokenHardwareID")
	defer span.End()

	if err := s.checkHardwareIDBanned(ctx, tx, tokenHardwareID); err != nil {
		return err
	}

	if storedHardwareID == nil {
		return apperrors.ErrHardwareIDIsInvalid
	}

	if *storedHardwareID != tokenHardwareID {
		return apperrors.ErrTokenHardwareIDIsInvalid
	}

	return nil
}

// verifyAndUpdateHardwareID validates HardwareID and updates it if necessary.
func (s *AuthenticationService) verifyAndUpdateHardwareID(
	ctx context.Context,
	tx *ent.Tx,
	user *dto.UserFullDTO,
	hardwareID string,
) error {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.verifyAndUpdateHardwareID")
	defer span.End()

	hardwareIDBanned, _ := s.bannedHardwareIDRepo.TxFindByHardwareID(ctx, tx, hardwareID)

	if hardwareIDBanned != nil {
		return apperrors.ErrHardwareIDBanned(hardwareIDBanned.BanReason)
	}

	if user.HardwareID == nil {
		encodedHardwareID := s.encodeHardwareID(ctx, hardwareID)
		user.HardwareID = &encodedHardwareID

		return s.authRepo.TxUpdateHardwareIDByID(ctx, tx, user.ID, encodedHardwareID)
	}

	// Verify hardware ID
	if !s.credentialsHelper.VerifyHardwareID(hardwareID, *user.HardwareID) {
		return apperrors.ErrUserWrongHardwareID
	}

	return nil
}

func (s *AuthenticationService) checkHardwareIDBanned(
	ctx context.Context,
	tx *ent.Tx,
	encryptedHardwareID string,
) error {
	originalHWID, err := s.credentialsHelper.DecodeHardwareID(encryptedHardwareID)
	if err != nil {
		return apperrors.ErrHardwareIDIsInvalid
	}

	hardwareIDBanned, _ := s.bannedHardwareIDRepo.TxFindByHardwareID(ctx, tx, originalHWID)

	if hardwareIDBanned != nil {
		return apperrors.ErrHardwareIDBanned(hardwareIDBanned.BanReason)
	}

	return nil
}

// isAccountLocked checks if user account is locked.
func (s *AuthenticationService) isAccountLocked(user *dto.UserDTO) bool {
	return user.AccountBlockedUntil != nil && user.AccountBlockedUntil.After(time.Now())
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
	user *dto.UserFullDTO,
) *domainservice.AuthenticationResult {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationService.createAuthResult")
	defer span.End()

	var token string

	if user.HardwareID != nil { // HardwareID must be updated if nil
		tokenData := &entity.TokenData{
			ID:         user.ID,
			Username:   user.Username,
			HardwareID: *user.HardwareID,
		}

		token = s.generateToken(ctx, tokenData)
	}

	online := s.getOnlineCount(ctx)

	return domainservice.NewAuthenticationResult(token, user, online)
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
