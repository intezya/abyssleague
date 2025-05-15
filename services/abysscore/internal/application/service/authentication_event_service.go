package applicationservice

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/userentity"
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	"github.com/intezya/abyssleague/services/abysscore/pkg/timeutils"
	"github.com/intezya/pkglib/logger"
	"golang.org/x/sync/errgroup"
	"time"
)

type AuthenticationEventService struct {
	userRepository repositoryports.UserRepository
}

func NewAuthenticationEventService(userRepository repositoryports.UserRepository) *AuthenticationEventService {
	return &AuthenticationEventService{userRepository: userRepository}
}

func (s *AuthenticationEventService) HandleRegistration(ctx context.Context, user *dto.UserDTO) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationEventService.HandleRegistration")
	defer span.End()

}

func (s *AuthenticationEventService) HandleLogin(ctx context.Context, user *dto.UserDTO) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationEventService.HandleLogin")
	defer span.End()

	tx, err := s.userRepository.WithTx(ctx)
	if err != nil {
		logger.Log.Warnln("failed to start transaction:", err)
		return
	}

	err = persistence.WithTx(
		ctx, tx, func(tx *ent.Tx) error {
			group, _ := errgroup.WithContext(ctx)

			group.Go(
				func() error {
					return s.processLoginStreakAndRewards(ctx, tx, user)
				},
			)
			group.Go(
				func() error {
					return s.processBanDecrementAfterLogin(ctx, tx, user)
				},
			)
			// TODO: Implement logic to add bonuses for user based on login streak

			return group.Wait()
		},
	)

	if err != nil {
		logger.Log.Errorln("error in authentication handler:", err)
	}
}

func (s *AuthenticationEventService) processLoginStreakAndRewards(
	ctx context.Context,
	tx *ent.Tx,
	user *dto.UserDTO,
) error {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationEventsHandlers.processLoginStreakAndRewards")
	defer span.End()

	if !timeutils.IsDayBeforeToday(user.LoginAt) {
		return nil
	}

	user.LoginStreak++
	user.LoginAt = time.Now()

	err := s.userRepository.TxUpdateLoginStreakLoginAtByID(ctx, tx, user.ID, user.LoginStreak, user.LoginAt)
	if err != nil {
		logger.Log.Warnw("failed to update login streak", "error", err, "userID", user.ID)
		return err
	}

	return nil
}

// processBanDecrementAfterLogin decrements ban levels if needed.
func (s *AuthenticationEventService) processBanDecrementAfterLogin(
	ctx context.Context,
	tx *ent.Tx,
	user *dto.UserDTO,
) error {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationEventsHandlers.processBanDecrementAfterLogin")
	defer span.End()

	now := time.Now()

	needsUpdate := false

	if user.AccountBlockedUntil != nil &&
		user.AccountBlockedUntil.Add(userentity.AccountBlockDecrementTime).Before(now) {
		if user.AccountBlockedLevel > 0 {
			user.AccountBlockedLevel--
		}

		user.AccountBlockedUntil = nil
		user.AccountBlockReason = nil
		needsUpdate = true
	}

	if user.SearchBlockedUntil != nil &&
		user.SearchBlockedUntil.Add(userentity.SearchBlockDecrementTime).Before(now) {
		if user.SearchBlockedLevel > 0 {
			user.SearchBlockedLevel--
		}

		user.SearchBlockedUntil = nil
		user.SearchBlockReason = nil
		needsUpdate = true
	}

	if needsUpdate {
		err := s.userRepository.TxSetBlockUntilAndLevelAndReasonFromUser(ctx, tx, user)
		if err != nil {
			logger.Log.Errorw("Failed to update block settings for user", "error", err, "userID", user.ID)
			return err
		}
	}

	return nil
}
