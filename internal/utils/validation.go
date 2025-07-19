package utils

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using the validator tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// GetValidationErrors extracts validation errors and formats them
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = getErrorMessage(e)
		}
	}
	
	return errors
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email address"
	case "min":
		return "This field must be at least " + e.Param() + " characters long"
	case "max":
		return "This field must be at most " + e.Param() + " characters long"
	case "len":
		return "This field must be exactly " + e.Param() + " characters long"
	case "oneof":
		return "This field must be one of: " + e.Param()
	case "url":
		return "Please enter a valid URL"
	case "e164":
		return "Please enter a valid phone number (with country code)"
	case "gtfield":
		return "This field must be greater than " + e.Param()
	case "gte":
		return "This field must be greater than or equal to " + e.Param()
	case "lte":
		return "This field must be less than or equal to " + e.Param()
	default:
		return "Invalid value"
	}
}

// BindAndValidate binds request data and validates it
func BindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	
	if err := ValidateStruct(req); err != nil {
		validationErrors := GetValidationErrors(err)
		return ValidationError(c, validationErrors)
	}
	
	return nil
}

// GenerateRandomToken generates a random token of specified length
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
