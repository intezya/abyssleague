package applicationservice

import (
	applicationerror "abysscore/internal/common/errors/application"
	"abysscore/internal/domain/dto"
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

type AuthenticationService struct {
	userRepository       repositoryports.UserRepository
	mainWebsocketService drivenports.WebsocketService
	credentialsHelper    domainservice.CredentialsHelper
	tokenHelper          domainservice.TokenHelper
}

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

func (a *AuthenticationService) Register(
	ctx context.Context,
	credentials *userentity.CredentialsDTO,
) (
	*domainservice.AuthenticationResult,
	error,
) {
	credentials.Password = tracer.TraceValue(ctx, "credentialsHelper.EncodePassword", func(ctx context.Context) string {
		return a.credentialsHelper.EncodePassword(credentials.Password)
	})

	credentials.Hwid = tracer.TraceValue(ctx, "credentialsHelper.EncodeHardwareID", func(ctx context.Context) string {
		return a.credentialsHelper.EncodeHardwareID(credentials.Hwid)
	})

	user, err := tracer.TraceFnWithResult(ctx, "userRepository.Create", func(ctx context.Context) (*userentity.AuthenticationData, error) {
		return a.userRepository.Create(credentials)
	})

	if err != nil {
		// Could be a username conflict or hardware ID conflict
		return nil, err
	}

	token := tracer.TraceValue(ctx, "tokenHelper.TokenGenerator", func(ctx context.Context) string {
		return a.tokenHelper.TokenGenerator(user.TokenData())
	})

	online := tracer.TraceValue(ctx, "mainWebsocketService.GetOnlineSoft", func(ctx context.Context) int {
		res := a.mainWebsocketService.GetOnlineSoft(ctx)
		return int(res.Online)
	})

	return domainservice.NewAuthenticationResult(token, nil, online), nil
}

func (a *AuthenticationService) Authenticate(ctx context.Context, credentials *userentity.CredentialsDTO) (
	*domainservice.AuthenticationResult,
	error,
) {
	lowerUsername := strings.ToLower(credentials.Username)

	authentication, err := tracer.TraceFnWithResult(ctx, "userRepository.FindAuthenticationByLowerUsername", func(ctx context.Context) (*userentity.AuthenticationData, error) {
		return a.userRepository.FindAuthenticationByLowerUsername(lowerUsername)
	})

	if err != nil {
		// User with the provided username was not found
		return nil, err
	}

	passwordOk := tracer.TraceValue(ctx, "authentication.ComparePassword", func(ctx context.Context) bool {
		return authentication.ComparePassword(credentials.Password, a.credentialsHelper.VerifyPassword)
	})

	if !passwordOk {
		return nil, applicationerror.ErrWrongPassword
	}

	ok, needsUpdate := tracer.TraceValueValue(ctx, "authentication.CompareHWID", func(ctx context.Context) (bool, bool) {
		return authentication.CompareHWID(credentials.Hwid, a.credentialsHelper.VerifyHardwareID)
	})

	if !ok {
		return nil, applicationerror.ErrUserWrongHwid
	}

	if needsUpdate {
		newHwid := tracer.TraceValue(ctx, "credentialsHelper.EncodeHardwareID", func(ctx context.Context) string {
			return a.credentialsHelper.EncodeHardwareID(credentials.Hwid)
		})

		authentication.SetHWID(newHwid)

		err = tracer.TraceValue(ctx, "userRepository.UpdateHWIDByID", func(ctx context.Context) error {
			return a.userRepository.UpdateHWIDByID(authentication.UserID(), newHwid)
		})

		if err != nil {
			return nil, err // hwid conflict
		}
	}

	token := tracer.TraceValue(ctx, "authentication.TokenGenerator", func(ctx context.Context) string {
		return a.tokenHelper.TokenGenerator(authentication.TokenData())
	})

	user, err := tracer.TraceFnWithResult(ctx, "userRepository.FindFullDTOById", func(ctx context.Context) (*dto.UserFullDTO, error) {
		return a.userRepository.FindFullDTOById(authentication.UserID())
	}) // Cannot return err cause it handled above

	if err != nil {
		logger.Log.Warnw("Failed to retrieve full user data", "error", err, "userID", authentication.UserID())
		// Continue with authentication even if user data retrieval fails
		// The token is still valid and can be used for authentication
	}

	go a.postLoginProcessing(ctx, user.UserDTO)

	online := tracer.TraceValue(ctx, "mainWebsocketService.GetOnlineSoft", func(ctx context.Context) int {
		res := a.mainWebsocketService.GetOnlineSoft(ctx)
		return int(res.Online)
	})

	return domainservice.NewAuthenticationResult(token, user, online), nil
}

func (a *AuthenticationService) ValidateToken(ctx context.Context, token string) (*dto.UserDTO, error) {
	data, err := tracer.TraceFnWithResult(ctx, "tokenHelper.ValidateToken", func(ctx context.Context) (*userentity.TokenData, error) {
		return a.tokenHelper.ValidateToken(token)
	})

	if err != nil {
		return nil, err
	}

	logger.Log.Debugw("authentication data received from token", "data", data)

	lowerUsername := strings.ToLower(data.Username)

	authentication, err := tracer.TraceFnWithResult(ctx, "userRepository.FindAuthenticationByLowerUsername", func(ctx context.Context) (*userentity.AuthenticationData, error) {
		return a.userRepository.FindAuthenticationByLowerUsername(lowerUsername)
	})

	if err != nil {
		// User from token not found in the database
		return nil, err
	}

	ok, needsUpdate := tracer.TraceValueValue(ctx, "authentication.CompareHWID", func(ctx context.Context) (bool, bool) {
		return authentication.CompareHWID(data.Hwid, a.credentialsHelper.VerifyHardwareID)
	})

	if !ok || needsUpdate {
		return nil, applicationerror.ErrTokenHwidIsInvalid
	}

	user, err := tracer.TraceFnWithResult(ctx, "userRepository.FindDTOById", func(ctx context.Context) (*dto.UserDTO, error) {
		return a.userRepository.FindDTOById(authentication.UserID())
	})

	if err != nil {
		logger.Log.Warnw("Failed to retrieve user data during token validation", "error", err, "userID", authentication.UserID())
		return nil, err
	}

	return user, nil
}

func (a *AuthenticationService) postLoginProcessing(ctx context.Context, user *dto.UserDTO) {
	if !timeutils.IsDayBeforeToday(user.LoginAt) {
		return
	}

	user.LoginStreak++
	user.LoginAt = time.Now()

	// TODO: here will be added bonuses for user (maybe some money, xp in stats, or badge if login streak too high)

	err := tracer.TraceValue(ctx, "userRepository.SetLoginStreakLoginAtByID", func(ctx context.Context) error {
		return a.userRepository.SetLoginStreakLoginAtByID(user.ID, user.LoginStreak, user.LoginAt)
	})

	if err != nil {
		logger.Log.Warnf("Unexpected error occurred: %v", err)
	}
}
