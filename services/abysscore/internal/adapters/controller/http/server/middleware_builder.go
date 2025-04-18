package server

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

// buildMiddleware creates a chain of middleware for the given route
func (m *MiddlewareLinker) buildMiddleware(route *Route) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create middleware chain based on route requirements
		var middlewareChain []fiber.Handler

		// Core middleware always applied
		middlewareChain = append(middlewareChain, m.loggingMiddleware.Handle())
		middlewareChain = append(middlewareChain, m.recoverMiddleware.Handle())

		// Add rate limiting based on route configuration
		switch route.RateLimit {
		case DefaultRateLimit:
			middlewareChain = append(middlewareChain, m.rateLimitMiddleware.HandleDefault())
		case AuthRateLimit:
			middlewareChain = append(middlewareChain, m.rateLimitMiddleware.HandleForAuth())
		default:
		}

		// Add authentication and authorization checks if required
		if route.RequireAuthentication {
			// Authentication check
			middlewareChain = append(middlewareChain, m.authenticationMiddleware.Handle())

			// Access level check if specified
			if route.AccessLevel != nil {
				middlewareChain = append(middlewareChain, createAccessLevelChecker(route.AccessLevel))
			}

			// Match requirement check if specified
			if route.MatchRequirement != MatchIrrelevant {
				middlewareChain = append(middlewareChain, createMatchRequirementChecker(route.MatchRequirement))
			}
		}

		// Finally add the route handler
		middlewareChain = append(middlewareChain, route.Handler)

		// Execute the middleware chain
		return executeMiddlewareChain(middlewareChain, c)
	}
}

// executeMiddlewareChain runs each middleware in sequence
func executeMiddlewareChain(chain []fiber.Handler, c *fiber.Ctx) error {
	for _, middleware := range chain {
		if err := middleware(c); err != nil {
			return err
		}
	}
	return nil
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
