// common/errors/base/error.go
package base

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"runtime"
	"time"
)

// Error represents a structured application error with stack trace capabilities
type Error struct {
	Custom    error                  // Public-facing error message
	Wrapped   error                  // Original error for debugging
	code      int                    // HTTP status code
	stack     []string               // Stack trace
	rawStack  string                 // Raw stack trace (for panics)
	timestamp time.Time              // When error occurred
	errorID   string                 // Unique error identifier
	metadata  map[string]interface{} // Additional error context
}

// ErrorResponse defines the structure of error JSON responses
type ErrorResponse struct {
	Message   string                 `json:"message"`
	Detail    string                 `json:"detail,omitempty"`
	Code      int                    `json:"code"`
	Path      string                 `json:"path"`
	Timestamp time.Time              `json:"timestamp"`
	ErrorID   string                 `json:"error_id"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Custom.Error()
}

// WithCode allows changing the status code of an existing error
func (e *Error) WithCode(code int) *Error {
	e.code = code
	return e
}

// SetErrorID sets the unique error identifier
func (e *Error) SetErrorID(id string) {
	e.errorID = id
}

// SetTimestamp sets the error timestamp
func (e *Error) SetTimestamp(t time.Time) {
	e.timestamp = t
}

// SetStackTrace sets the raw stack trace (useful for panics)
func (e *Error) SetStackTrace(stack string) {
	e.rawStack = stack
}

// SetMetadata sets additional context information for the error
func (e *Error) SetMetadata(data map[string]interface{}) {
	e.metadata = data
}

// AddMetadata adds a single key-value pair to the error metadata
func (e *Error) AddMetadata(key string, value interface{}) *Error {
	if e.metadata == nil {
		e.metadata = make(map[string]interface{})
	}
	e.metadata[key] = value
	return e
}

// NewError creates a structured error with stack trace
func NewError(custom error, wrapped error, code int) *Error {
	// Generate stack trace
	stackBuf := make([]uintptr, 20)
	length := runtime.Callers(2, stackBuf)
	stack := make([]string, 0, length)

	frames := runtime.CallersFrames(stackBuf[:length])
	for {
		frame, more := frames.Next()
		// Skip system and library frames
		if !isSystemFrame(frame.Function) {
			stack = append(stack, fmt.Sprintf("%s:%d", frame.Function, frame.Line))
		}
		if !more {
			break
		}
	}

	// Generate unique error ID
	errorID := generateErrorID()

	return &Error{
		Custom:    custom,
		Wrapped:   wrapped,
		code:      code,
		stack:     stack,
		timestamp: time.Now(),
		errorID:   errorID,
		metadata:  make(map[string]interface{}),
	}
}

// StatusCode returns the HTTP status code
func (e *Error) StatusCode() int {
	return e.code
}

// Message returns the user-facing error message
func (e *Error) Message() string {
	if e.Custom != nil {
		return e.Custom.Error()
	}
	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}
	return "unknown error"
}

// Detail returns detailed technical error information
func (e *Error) Detail() string {
	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}
	if e.Custom != nil {
		return e.Custom.Error()
	}
	return ""
}

// ErrorID returns the unique error identifier
func (e *Error) ErrorID() string {
	return e.errorID
}

// Timestamp returns when the error occurred
func (e *Error) Timestamp() time.Time {
	return e.timestamp
}

// StackTrace returns the error stack trace
func (e *Error) StackTrace() []string {
	return e.stack
}

// RawStackTrace returns the raw stack trace string (for panics)
func (e *Error) RawStackTrace() string {
	return e.rawStack
}

// Metadata returns additional error context
func (e *Error) Metadata() map[string]interface{} {
	return e.metadata
}

// ToErrorResponse converts the error to HTTP response
func (e *Error) ToErrorResponse(c *fiber.Ctx) error {
	response := &ErrorResponse{
		Message:   e.Message(),
		Detail:    e.Detail(),
		Code:      e.StatusCode(),
		Path:      c.Path(),
		Timestamp: e.timestamp,
		ErrorID:   e.errorID,
	}

	// Only include metadata in response if it exists and is not empty
	if len(e.metadata) > 0 {
		response.Metadata = e.metadata
	}

	return c.Status(e.StatusCode()).JSON(response)
}

// Helper functions
func isSystemFrame(funcName string) bool {
	// Skip standard library and runtime frames
	return funcName == "runtime.goexit" ||
		funcName == "runtime.main" ||
		funcName == "runtime.gopanic"
}

func generateErrorID() string {
	return fmt.Sprintf("err_%d_%s", time.Now().UnixNano(), uuid.New())
}
