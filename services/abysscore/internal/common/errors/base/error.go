package base

import (
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Custom  error
	Wrapped error
	code    int
}

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"`
	Path    string `json:"path"`
}

func (e *Error) Error() string {
	return e.Custom.Error()
}

func NewError(custom error, wrapped error, code int) *Error {
	return &Error{Custom: custom, Wrapped: wrapped, code: code}
}

func (e *Error) StatusCode() int {
	return e.code
}

func (e *Error) Message() string {
	if e.Custom != nil {
		return e.Custom.Error()
	}

	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}

	return ""
}

func (e *Error) Detail() string {
	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}

	if e.Custom != nil {
		return e.Custom.Error()
	}

	return ""
}

func (e *Error) ToErrorResponse(c *fiber.Ctx) error {
	return c.Status(e.StatusCode()).JSON(&ErrorResponse{
		Message: e.Message(),
		Detail:  e.Detail(),
		Code:    e.StatusCode(),
		Path:    c.Path(),
	})
}
