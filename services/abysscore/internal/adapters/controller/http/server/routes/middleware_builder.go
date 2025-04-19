package routes

import (
	"abysscore/internal/adapters/controller/http/middleware"
	adaptererror "abysscore/internal/common/errors/adapter"
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent/schema/access_level"
	"github.com/gofiber/fiber/v2"
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

// createAccessLevelChecker creates a middleware to check user access level
func createAccessLevelChecker(requiredLevel *access_level.AccessLevel) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.UserContext().Value(middleware.UserCtxKey).(*dto.UserDTO)

		if user.AccessLevel < *requiredLevel {
			return adaptererror.InsufficientAccessLevel.ToErrorResponse(c)
		}

		return c.Next()
	}
}

// createMatchRequirementChecker creates a middleware to check user match state
func createMatchRequirementChecker(matchRequirement MatchRequirement) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.UserContext().Value(middleware.UserCtxKey).(*dto.UserDTO)

		switch matchRequirement {
		case MustBeInMatch:
			if user.CurrentMatchID == nil {
				return adaptererror.UserMatchStateError(
					fiber.NewError(fiber.StatusForbidden, "user must be in match"),
				).ToErrorResponse(c)
			}
		case MustNotBeInMatch:
			if user.CurrentMatchID != nil {
				return adaptererror.UserMatchStateError(
					fiber.NewError(fiber.StatusForbidden, "user must not be in match"),
				).ToErrorResponse(c)
			}
		default:
		}

		return c.Next()
	}
}
