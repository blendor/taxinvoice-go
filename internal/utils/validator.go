package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation for state code
	_ = validate.RegisterValidation("statecode", validateStateCode)
}

// ValidateStruct validates a struct using tags
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return NewInternalServerError("Invalid validation error", err)
		}

		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, formatValidationError(err))
		}

		return NewBadRequestError(strings.Join(errorMessages, "; "))
	}
	return nil
}

// formatValidationError formats a single validation error into a human-readable message
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, err.Param())
	case "statecode":
		return fmt.Sprintf("%s must be a valid two-letter state code", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, err.Tag())
	}
}

// validateStateCode is a custom validator for two-letter state codes
func validateStateCode(fl validator.FieldLevel) bool {
	stateCode := fl.Field().String()
	match, _ := regexp.MatchString(`^[A-Z]{2}$`, stateCode)
	return match
}

// ValidateField validates a single field
func ValidateField(field interface{}, tag string) error {
	err := validate.Var(field, tag)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return NewInternalServerError("Invalid validation error", err)
		}
		return NewBadRequestError(formatValidationError(err.(validator.ValidationErrors)[0]))
	}
	return nil
}