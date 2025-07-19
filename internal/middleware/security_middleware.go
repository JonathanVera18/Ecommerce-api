package middleware

import (
	"github.com/labstack/echo/v4"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// X-Content-Type-Options: Prevent MIME type sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// X-Frame-Options: Prevent clickjacking
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// X-XSS-Protection: Enable XSS filtering
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Referrer-Policy: Control referrer information
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Content-Security-Policy: Prevent code injection
			csp := "default-src 'self'; " +
				"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
				"style-src 'self' 'unsafe-inline'; " +
				"img-src 'self' data: https:; " +
				"font-src 'self' data:; " +
				"connect-src 'self'; " +
				"frame-ancestors 'none';"
			c.Response().Header().Set("Content-Security-Policy", csp)

			// Strict-Transport-Security: Enforce HTTPS (only if HTTPS)
			if c.Request().TLS != nil || c.Request().Header.Get("X-Forwarded-Proto") == "https" {
				c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}

			// Remove server information
			c.Response().Header().Set("Server", "")

			// Permissions-Policy: Control browser features
			c.Response().Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

			return next(c)
		}
	}
}

// HTTPS redirect middleware
func HTTPSRedirect() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip redirect for health checks and development
			if c.Request().URL.Path == "/health" || c.Request().Header.Get("X-Forwarded-Proto") == "https" {
				return next(c)
			}

			// Only redirect in production
			if c.Request().Header.Get("X-Forwarded-Proto") == "http" {
				return c.Redirect(301, "https://"+c.Request().Host+c.Request().RequestURI)
			}

			return next(c)
		}
	}
}
