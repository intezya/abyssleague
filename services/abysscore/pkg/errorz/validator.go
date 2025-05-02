package errorz

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// ValidatorInterface allows us to mock the validator for testing.
type ValidatorInterface interface {
	Struct(s interface{}) error
}

// validatorInstance is the singleton validator instance.
var validatorInstance ValidatorInterface

// SetValidator sets the validator instance to be used.
func SetValidator(v ValidatorInterface) {
	validatorInstance = v
}

// GetValidator returns the validator instance.
func GetValidator() ValidatorInterface {
	if validatorInstance == nil {
		// Default to standard validator if not set
		validatorInstance = validator.New()
	}

	return validatorInstance
}

// ValidateStruct validates a struct against validation tags.
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}

// ValidateJSON validates a data transfer object and returns a validation error if invalid.
func ValidateJSON(dto interface{}) error {
	if err := ValidateStruct(dto); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorMessages := make([]string, len(validationErrors))
			for i, err := range validationErrors {
				errorMessages[i] = formatValidationError(err)
			}

			// Create and return a validation error
			return validationError(err).WithValidations(errorMessages)
		}

		// Generic validation error
		return validationError(err)
	}

	return nil
}

// formatValidationError formats a validation error in a user-friendly way.
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
