package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/base"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/errorutils"
	"github.com/intezya/pkglib/logger"
)

type RecoverMiddleware struct {
	requestIdConfig requestid.Config
}

func NewRecoverMiddleware(requestIdConfig requestid.Config) *RecoverMiddleware {
	return &RecoverMiddleware{
		requestIdConfig: requestIdConfig,
	}
}

func (r *RecoverMiddleware) Handle() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				stackTrace := debug.Stack()
				requestID := c.Locals(r.requestIdConfig.ContextKey)
				timestamp := time.Now()
				errorID := fmt.Sprintf("panic_%s_%d", requestID, timestamp.UnixNano())

				// Convert panic to error
				var recErr error
				switch v := rec.(type) {
				case error:
					recErr = v
				default:
					recErr = fmt.Errorf("%v", rec)
				}

				// Create a structured error with additional context
				structuredErr := base.NewError(
					fmt.Errorf("server panic: %w", recErr),
					recErr,
					fiber.StatusInternalServerError,
				)

				// Add metadata to the error
				structuredErr.SetErrorID(errorID)
				structuredErr.SetTimestamp(timestamp)
				structuredErr.SetStackTrace(string(stackTrace))
				structuredErr.SetMetadata(map[string]interface{}{
					"request_id": requestID,
					"url":        c.OriginalURL(),
					"method":     c.Method(),
					"ip":         c.IP(),
					"user_agent": c.Get("User-Agent"),
				})

				// Log the error with full context
				logger.Log.With(
					"error_id", errorID,
					"error", fmt.Sprintf("%v", recErr),
					"stacktrace", string(stackTrace),
					"url", c.OriginalURL(),
					"method", c.Method(),
					"ip", c.IP(),
					"request_id", requestID,
					"timestamp", timestamp,
				).Error("panic recovered")

				// Use the utils package to log structured error
				errorutils.LogError(structuredErr)

				// Return error response
				err = structuredErr.ToErrorResponse(c)
			}
		}()

		// Call the next handler
		return c.Next()
	}
}
