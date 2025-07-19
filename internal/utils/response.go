package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

// SuccessResponse sends a successful JSON response
func SuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessResponseWithMeta sends a successful JSON response with metadata
func SuccessResponseWithMeta(c echo.Context, message string, data interface{}, meta interface{}) error {
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// CreatedResponse sends a created JSON response
func CreatedResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, models.ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// BadRequestError sends a bad request error response
func BadRequestError(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusBadRequest, message)
}

// UnauthorizedError sends an unauthorized error response
func UnauthorizedError(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenError sends a forbidden error response
func ForbiddenError(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusForbidden, message)
}

// NotFoundError sends a not found error response
func NotFoundError(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusNotFound, message)
}

// ConflictError sends a conflict error response
func ConflictError(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusConflict, message)
}

// InternalServerError sends an internal server error response
func InternalServerError(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusInternalServerError, message)
}

// ValidationError sends a validation error response with details
func ValidationError(c echo.Context, errors map[string]string) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   "Validation failed",
		"details": errors,
	})
}
