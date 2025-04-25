package base

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/pkglib/logger"
	"net/http"
	"time"
)

func ParseErrorOrInternalResponse(err error, c *fiber.Ctx) error {
	var custom *Error

	var customValidation *ValidationError

	if !errors.As(err, &custom) && !errors.As(err, &customValidation) {
		logger.Log.Warnw(
			"returned error not recognized",
			"err", err,
			"err_message", err.Error(),
			"path", c.Path(),
		)

		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponse{
			Message:   err.Error(),
			Detail:    "error not recognized",
			Code:      http.StatusInternalServerError,
			Path:      c.Path(),
			Timestamp: time.Now(),
			ErrorID:   generateErrorID(),
			Metadata:  nil,
		})
	}

	if custom != nil {
		return custom.ToErrorResponse(c)
	}

	if customValidation != nil {
		return customValidation.ToErrorResponse(c)
	}

	panic("unreachable")
}
