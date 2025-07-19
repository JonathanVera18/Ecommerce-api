package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// RateLimitConfig holds the configuration for rate limiting
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
	SkipSuccessful    bool
}

// DefaultRateLimitConfig returns a default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 100,
		BurstSize:         10,
		SkipSuccessful:    false,
	}
}

// RateLimit returns a rate limiting middleware
func RateLimit() echo.MiddlewareFunc {
	return RateLimitWithConfig(DefaultRateLimitConfig())
}

// RateLimitWithConfig returns a rate limiting middleware with custom configuration
func RateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	// Create a rate limiter per IP address
	limiters := make(map[string]*rate.Limiter)

	// Convert requests per minute to requests per second
	requestsPerSecond := float64(config.RequestsPerMinute) / 60.0

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get client IP
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			// Get or create rate limiter for this IP
			limiter, exists := limiters[clientIP]
			if !exists {
				limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), config.BurstSize)
				limiters[clientIP] = limiter
			}

			// Check if request is allowed
			if !limiter.Allow() {
				// Add rate limit headers
				c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
				c.Response().Header().Set("X-RateLimit-Remaining", "0")
				c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}

			// Add rate limit headers for successful requests
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))

			return next(c)
		}
	}
}

// AuthRateLimit returns a stricter rate limit for authentication endpoints
func AuthRateLimit() echo.MiddlewareFunc {
	return RateLimitWithConfig(RateLimitConfig{
		RequestsPerMinute: 30, // More restrictive for auth endpoints
		BurstSize:         5,
		SkipSuccessful:    false,
	})
}

// APIRateLimit returns a general rate limit for API endpoints
func APIRateLimit() echo.MiddlewareFunc {
	return RateLimitWithConfig(RateLimitConfig{
		RequestsPerMinute: 100,
		BurstSize:         20,
		SkipSuccessful:    false,
	})
}
