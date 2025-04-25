package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	adaptererror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/adapter"
	"github.com/intezya/abyssleague/services/abysscore/internal/common/errors/base"
	"github.com/intezya/pkglib/logger"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

var v = NewValidator() //nolint:varnamelen // disabled because it is pkg (truth-source) code

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

var errValidationFailed = adaptererror.UnprocessableEntity(errors.New("validation error"))

func ValidateJSON(dto interface{}) error {
	if err := v.ValidateStruct(dto); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorMessages := make([]string, len(validationErrors))
			for i, err := range validationErrors {
				errorMessages[i] = formatValidationError(err)
			}

			logger.Log.Debugw("many validation errors", "errors", errorMessages, "err", err)

			return base.NewValidationError(err, errorMessages)
		}

		logger.Log.Debugw("validation error", "err", err)

		return errValidationFailed
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
