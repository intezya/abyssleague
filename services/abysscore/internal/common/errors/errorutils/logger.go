package errorutils

import (
	"errors"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/base"
	"github.com/intezya/pkglib/logger"
)

// LogError logs an error with its details and stack trace.
func LogError(err error) {
	// Try to cast to our custom error type
	var customErr *base.Error
	if errors.As(err, &customErr) {
		// Log structured error with stack trace
		logger.Log.Errorw(
			customErr.Message(),
			"error_id", customErr.ErrorID(),
			"details", customErr.Detail(),
			"stack_trace", customErr.StackTrace(),
			"code", customErr.StatusCode,
			"timestamp", customErr.Timestamp,
		)
	} else {
		// Log standard error
		logger.Log.Errorw(
			"Unstructured error",
			"error", err.Error(),
		)
	}
}
