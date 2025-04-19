package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is a wrapper around the validator.Validate
type Validator struct {
	validate *validator.Validate
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error returns the error message
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	
	var sb strings.Builder
	sb.WriteString("validation failed: ")
	
	for i, err := range ve {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	
	return sb.String()
}

// New creates a new validator
func New() *Validator {
	validate := validator.New()
	
	// Register custom validation functions here
	
	// Use JSON tag names for validation errors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
	
	return &Validator{
		validate: validate,
	}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}
	
	var validationErrors ValidationErrors
	
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Message: getErrorMessage(err),
		})
	}
	
	return validationErrors
}

// getErrorMessage returns a human-readable error message for a validation error
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long", err.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", err.Param())
	default:
		return fmt.Sprintf("Failed validation on %s", err.Tag())
	}
}
