package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtService *utils.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Success: false,
					Error:   "Authorization header required",
				})
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Success: false,
					Error:   "Invalid authorization header format",
				})
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Success: false,
					Error:   "Token required",
				})
			}

			// Validate token
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Success: false,
					Error:   "Invalid or expired token",
				})
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

// JWTAuth is an alias for AuthMiddleware for compatibility
func JWTAuth(jwtService *utils.JWTService) echo.MiddlewareFunc {
	return AuthMiddleware(jwtService)
}

// OptionalAuthMiddleware validates JWT tokens if present, but doesn't require them
func OptionalAuthMiddleware(jwtService *utils.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return next(c)
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return next(c)
			}

			// Validate token
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				return next(c)
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}
