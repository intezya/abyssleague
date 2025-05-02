package errorz

import (
	"net/http"
)

// Common error generators for different application layers.
var (
	// Repository Errors.

	// NotFound creates a generic "not found" error.
	NotFound = func(entity string, detail error) *Error {
		return New(
			entity+" not found",
			detail,
			ErrorTypeRepository,
			http.StatusNotFound,
		)
	}

	// Conflict creates a generic conflict error.
	Conflict = func(message string, detail error) *Error {
		return New(
			message,
			detail,
			ErrorTypeRepository,
			http.StatusConflict,
		)
	}

	// Adapter (HTTP/API) Errors.

	// TooManyRequests creates a general rate limit error.
	TooManyRequests = func(err error) *Error {
		return New(
			"too many requests",
			err,
			ErrorTypeApplication,
			http.StatusTooManyRequests,
		)
	}

	// InternalError creates an internal server error.
	InternalError = func(detail error) *Error {
		return New(
			"internal server error",
			detail,
			ErrorTypeInternal,
			http.StatusInternalServerError,
		)
	}

	// Forbidden creates a permission error.
	Forbidden = func(message string, detail error) *Error {
		return New(
			message,
			detail,
			ErrorTypeAuthorization,
			http.StatusForbidden,
		)
	}

	// BadRequest creates a bad request error.
	BadRequest = func(detail error) *Error {
		return New(
			"bad request",
			detail,
			ErrorTypeValidation,
			http.StatusBadRequest,
		)
	}

	// UnprocessableEntity creates an unprocessable entity error.
	UnprocessableEntity = func(detail error) *Error {
		return New(
			"unprocessable entity",
			detail,
			ErrorTypeValidation,
			http.StatusUnprocessableEntity,
		)
	}

	// Unauthorized creates an unauthorized error.
	Unauthorized = func(detail error) *Error {
		return New(
			"unauthorized",
			detail,
			ErrorTypeAuthorization,
			http.StatusUnauthorized,
		)
	}

	// ServiceUnavailable creates a service unavailable error.
	ServiceUnavailable = func(detail error) *Error {
		return New(
			"service unavailable",
			detail,
			ErrorTypeInternal,
			http.StatusServiceUnavailable,
		)
	}
)

var validationError = func(err error) *Error {
	return New(
		"validation error",
		err,
		ErrorTypeValidation,
		http.StatusBadRequest,
	)
}
