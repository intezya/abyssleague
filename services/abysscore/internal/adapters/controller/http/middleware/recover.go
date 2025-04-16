package middleware

import (
	adaptererror "abysscore/common/errors/adapter"
	"abysscore/common/errors/base"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/intezya/pkglib/logger"
	"runtime/debug"
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

				recErr, ok := rec.(error)
				if !ok {
					recErr = fmt.Errorf("%v", rec)
				}

				requestID := c.Locals(r.requestIdConfig.ContextKey)

				logger.Log.With(
					"error", fmt.Sprintf("%v", recErr),
					"stacktrace", string(stackTrace),
					"url", c.OriginalURL(),
					"method", c.Method(),
					"ip", c.IP(),
					"request_id", requestID,
				).Error("panic recovered")

				// Set the error to be returned from the middleware
				err = base.ParseErrorOrInternalResponse(adaptererror.InternalServerError, c)
			}
		}()

		// Call the next handler
		return c.Next()
	}
}
