package validator

import (
	adaptererror "abysscore/internal/common/errors/adapter"
	"abysscore/internal/common/errors/base"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

var v = NewValidator()

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

var validationError = adaptererror.UnprocessableEntity(errors.New("validation error"))

func ValidateJSON(dto interface{}, c *fiber.Ctx) error {
	if err := v.ValidateStruct(dto); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorMessages := make([]string, len(validationErrors))
			for i, err := range validationErrors {
				errorMessages[i] = formatValidationError(err)
			}

			return base.NewValidationError(nil, errorMessages)
		}
		return validationError
	}

	return nil
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " required"
	case "email":
		return err.Field() + " should be a valid email address"
	case "min":
		return err.Field() + " should be at least " + err.Param()
	case "max":
		return err.Field() + " should be less than " + err.Param()
	case "gte":
		return err.Field() + " should be greater than " + err.Param()
	case "lte":
		return err.Field() + " should be less than " + err.Param()
	default:
		return err.Field() + " validation error: " + err.Tag()
	}
}
