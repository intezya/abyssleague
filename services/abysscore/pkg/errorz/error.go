package errorz

import (
	"fmt"
	"runtime"
	"time"
)

// ErrorType represents the category of error.
type ErrorType string

// Define error categories.
const (
	// ErrorTypeValidation represents validation errors (bad input).
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeRepository represents database or storage errors.
	ErrorTypeRepository ErrorType = "repository"
	// ErrorTypeApplication represents business logic errors.
	ErrorTypeApplication ErrorType = "application"
	// ErrorTypeAuthorization represents authentication/authorization errors.
	ErrorTypeAuthorization ErrorType = "authorization"
	// ErrorTypeInternal represents unexpected system errors.
	ErrorTypeInternal ErrorType = "internal"
)

// Error is the central error type for the entire application.
type Error struct {
	Message     string                 // User-facing error message
	Detail      error                  // Original/technical error (for logs)
	StatusCode  int                    // HTTP status code
	ErrorType   ErrorType              // Category of error
	Stack       []string               // Stack trace
	RawStack    string                 // Raw stack trace (for panics)
	Timestamp   time.Time              // When error occurred
	ErrorID     string                 // Unique error identifier
	Metadata    map[string]interface{} // Additional error context
	Validations []string               // Validation-specific error details
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// New creates a new application error.
func New(message string, detail error, errorType ErrorType, code int) *Error {
	// Generate stack trace
	const stackDepth = 20
	stackBuf := make([]uintptr, stackDepth)
	length := runtime.Callers(2, stackBuf) //nolint:mnd // skip this and previous func calls
	stack := make([]string, 0, length)

	frames := runtime.CallersFrames(stackBuf[:length])

	for {
		frame, more := frames.Next()
		// Skip system frames
		if !isSystemFrame(frame.Function) {
			stack = append(stack, fmt.Sprintf("%s:%d", frame.Function, frame.Line))
		}

		if !more {
			break
		}
	}

	return &Error{
		Message:    message,
		Detail:     detail,
		StatusCode: code,
		ErrorType:  errorType,
		Stack:      stack,
		Timestamp:  time.Now(),
		ErrorID:    generateErrorID(),
		Metadata:   make(map[string]interface{}),
	}
}

func (e *Error) WithCode(code int) *Error {
	e.StatusCode = code

	return e
}

// WithMetadata adds contextual data to the error.
func (e *Error) WithMetadata(key string, value interface{}) *Error {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}

	e.Metadata[key] = value

	return e
}

// WithValidations adds validation errors.
func (e *Error) WithValidations(validations []string) *Error {
	e.Validations = validations

	return e
}

// ToResponse converts the error to an HTTP response.
func (e *Error) ToResponse(ctx Context) error {
	response := &ErrorResponse{
		Message:   e.Message,
		Code:      e.StatusCode,
		Path:      ctx.Path(),
		Timestamp: e.Timestamp,
		ErrorID:   e.ErrorID,
		Type:      string(e.ErrorType),
	}

	// Only include details for non-production environments
	if e.Detail != nil {
		response.Detail = e.Detail.Error()
	}

	// Add validation errors if present
	if len(e.Validations) > 0 {
		response.Validations = e.Validations
	}

	// Add metadata if present
	if len(e.Metadata) > 0 {
		response.Metadata = e.Metadata
	}

	return ctx.Status(e.StatusCode).JSON(response)
}
