package apperrors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/pkg/errorz"
	"github.com/intezya/pkglib/logger"
)

// FiberContext adapts Fiber's context to our errors package.
type FiberContext struct {
	ctx *fiber.Ctx
}

// Path returns the current request path.
func (c *FiberContext) Path() string {
	return c.ctx.Path()
}

// Status sets the HTTP status code.
func (c *FiberContext) Status(code int) errorz.Context {
	c.ctx.Status(code)

	return c
}

// JSON sends a JSON response.
func (c *FiberContext) JSON(data interface{}) error {
	return c.ctx.JSON(data)
}

// NewContext creates a new adapter for Fiber's context.
func NewContext(c *fiber.Ctx) errorz.Context {
	return &FiberContext{ctx: c}
}

// HandleError handles any error and sends an appropriate response.
func HandleError(err error, c *fiber.Ctx) error {
	return errorz.Handle(err, NewContext(c))
}

// LogError logs an error with all available details.
func LogError(err error) {
	var appErr *errorz.Error

	if errors.As(err, &appErr) {
		// Log structured error with stack trace
		logger.Log.Errorw(
			appErr.Message,
			"error_id", appErr.ErrorID,
			"details", appErr.Detail,
			"stack_trace", appErr.Stack,
			"code", appErr.StatusCode,
			"timestamp", appErr.Timestamp,
			"type", appErr.ErrorType,
		)
	} else {
		// Log standard error
		logger.Log.Errorw(
			"Unstructured error",
			"error", err,
		)
	}
}
