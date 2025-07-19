package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

// RequireRole creates middleware that requires specific user roles
func RequireRole(allowedRoles ...models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from context (set by auth middleware)
			userRole, ok := c.Get("user_role").(models.UserRole)
			if !ok {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Success: false,
					Error:   "Authentication required",
				})
			}

			// Check if user role is allowed
			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, models.ErrorResponse{
				Success: false,
				Error:   "Insufficient permissions",
			})
		}
	}
}

// RequireAdmin middleware that requires admin role
func RequireAdmin() echo.MiddlewareFunc {
	return RequireRole(models.RoleAdmin)
}

// RequireSeller middleware that requires seller role (or admin)
func RequireSeller() echo.MiddlewareFunc {
	return RequireRole(models.RoleSeller, models.RoleAdmin)
}

// RequireCustomer middleware that requires customer role (or admin)
func RequireCustomer() echo.MiddlewareFunc {
	return RequireRole(models.RoleCustomer, models.RoleAdmin)
}

// RequireSellerOrAdmin middleware that requires seller or admin role
func RequireSellerOrAdmin() echo.MiddlewareFunc {
	return RequireRole(models.RoleSeller, models.RoleAdmin)
}

// RequireCustomerOrAdmin middleware that requires customer or admin role
func RequireCustomerOrAdmin() echo.MiddlewareFunc {
	return RequireRole(models.RoleCustomer, models.RoleAdmin)
}
