package errorz

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func isSystemFrame(funcName string) bool {
	return funcName == "runtime.goexit" ||
		funcName == "runtime.main" ||
		funcName == "runtime.gopanic"
}

func generateErrorID() string {
	return fmt.Sprintf("err_%d_%s", time.Now().UnixNano(), uuid.New())
}

// Handle processes any error and converts it to the appropriate response
func Handle(err error, c Context) error {
	// If it's already our custom error type, just return it
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.ToResponse(c)
	}

	// Otherwise, wrap it as an internal server error
	return New(
		"Internal server error",
		err,
		ErrorTypeInternal,
		http.StatusInternalServerError,
	).ToResponse(c)
}
