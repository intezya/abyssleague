package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/middleware"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/schema/access_level"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
	"github.com/intezya/pkglib/logger"
)

var (
	errUserMustBeInMatch    = errors.New("user must be in match")
	errUserMustNotBeInMatch = errors.New("user must not be in match")
)

type MiddlewareLinker struct {
	loggingMiddleware        *middleware.LoggingMiddleware
	recoverMiddleware        *middleware.RecoverMiddleware
	rateLimitMiddleware      *middleware.RateLimitMiddleware
	authenticationMiddleware *middleware.AuthenticationMiddleware
}

func NewMiddlewareLinker(
	loggingMiddleware *middleware.LoggingMiddleware,
	recoverMiddleware *middleware.RecoverMiddleware,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
	authenticationMiddleware *middleware.AuthenticationMiddleware,
) *MiddlewareLinker {
	return &MiddlewareLinker{
		loggingMiddleware:        loggingMiddleware,
		recoverMiddleware:        recoverMiddleware,
		rateLimitMiddleware:      rateLimitMiddleware,
		authenticationMiddleware: authenticationMiddleware,
	}
}

// createAccessLevelChecker creates a middleware to check user access level.
func createAccessLevelChecker(requiredLevel *access_level.AccessLevel) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.UserContext().Value(middleware.UserCtxKey).(*dto.UserDTO)

		if !ok {
			logger.Log.Error("mismatched client type for middleware")

			return apperrors.HandleError(apperrors.InternalServerError, c)
		}

		if user.AccessLevel < *requiredLevel {
			return apperrors.HandleError(apperrors.ForbiddenByInsufficientAccessLevel, c)
		}

		return c.Next()
	}
}

// createMatchRequirementChecker creates a middleware to check user match state.
func createMatchRequirementChecker(matchRequirement MatchRequirement) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.UserContext().Value(middleware.UserCtxKey).(*dto.UserDTO)

		if !ok {
			logger.Log.Error("mismatched client type for middleware")

			return apperrors.HandleError(apperrors.InternalServerError, c)
		}

		switch matchRequirement {
		case MatchIrrelevant:
		case MustBeInMatch:
			if user.CurrentMatchID == nil {
				return apperrors.HandleError(apperrors.WrapUserMatchStateError(errUserMustBeInMatch), c)
			}
		case MustNotBeInMatch:
			if user.CurrentMatchID != nil {
				return apperrors.HandleError(apperrors.WrapUserMatchStateError(errUserMustNotBeInMatch), c)
			}
		default:
		}

		return c.Next()
	}
}
