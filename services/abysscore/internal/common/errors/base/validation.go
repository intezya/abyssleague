package base

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

const unprocessableEntityErrorMessage = "unprocessable entity"

var unprocessableEntityError = errors.New(unprocessableEntityErrorMessage)

type ValidationError struct {
	Custom  error
	Wrapped error
	errors  []string
	code    int
}

type ValidationErrorResponse struct {
	Message string   `json:"message"`
	Detail  string   `json:"detail"`
	Errors  []string `json:"errors"`
	Code    int      `json:"code"`
	Path    string   `json:"path"`
}

func NewValidationError(wrapped error, errors []string) error {
	return &ValidationError{
		Custom:  unprocessableEntityError,
		Wrapped: wrapped,
		errors:  errors,
		code:    fiber.StatusUnprocessableEntity,
	}
}

func (e *ValidationError) Error() string {
	return e.Custom.Error()
}

func (e *ValidationError) StatusCode() int {
	return e.code
}

func (e *ValidationError) Message() string {
	if e.Custom != nil {
		return e.Custom.Error()
	}
	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}
	return ""
}

func (e *ValidationError) Detail() string {
	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}
	if e.Custom != nil {
		return e.Custom.Error()
	}
	return ""
}

func (e *ValidationError) ToErrorResponse(c *fiber.Ctx) error {
	return c.Status(e.StatusCode()).JSON(&ValidationErrorResponse{
		Message: e.Message(),
		Detail:  e.Detail(),
		Errors:  e.errors,
		Code:    e.StatusCode(),
		Path:    c.Path(),
	})
}
