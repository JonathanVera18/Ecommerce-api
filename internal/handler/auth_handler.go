package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type authHandler struct {
	authService service.AuthService
}

// AuthHandler type alias for the concrete auth handler
type AuthHandler = authHandler

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *authHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	response, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			return utils.ConflictError(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to register user")
	}

	return utils.CreatedResponse(c, "User registered successfully", response)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		return utils.UnauthorizedError(c, err.Error())
	}

	return utils.SuccessResponse(c, "Login successful", response)
}

// RefreshToken handles JWT token refresh
// @Summary Refresh JWT token
// @Description Refresh an existing JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (h *authHandler) RefreshToken(c echo.Context) error {
	// Get token from header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return utils.UnauthorizedError(c, "Authorization header required")
	}

	token := authHeader[7:] // Remove "Bearer " prefix
	
	newToken, err := h.authService.RefreshToken(c.Request().Context(), token)
	if err != nil {
		return utils.UnauthorizedError(c, "Invalid or expired token")
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", map[string]string{
		"token": newToken,
	})
}

// Logout handles user logout
// @Summary User logout
// @Description Logout user (invalidate token)
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/logout [post]
func (h *authHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	
	err := h.authService.Logout(c.Request().Context(), userID)
	if err != nil {
		return utils.InternalServerError(c, "Failed to logout")
	}

	return utils.SuccessResponse(c, "Logout successful", nil)
}

// GetProfile handles getting current user profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/profile [get]
func (h *authHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	
	user, err := h.authService.GetCurrentUser(c.Request().Context(), userID)
	if err != nil {
		return utils.NotFoundError(c, "User not found")
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", user)
}

// ChangePassword handles password change
// @Summary Change user password
// @Description Change the password of the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.PasswordChangeRequest true "Password change request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/change-password [post]
func (h *authHandler) ChangePassword(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	
	var req models.PasswordChangeRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	err := h.authService.ChangePassword(c.Request().Context(), userID, &req)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			return utils.BadRequestError(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to change password")
	}

	return utils.SuccessResponse(c, "Password changed successfully", nil)
}

// ForgotPassword handles password reset requests
// @Summary Request password reset
// @Description Send a password reset link to the user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body struct { Email string `json:"email" validate:"required,email"` } true "Email address"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *authHandler) ForgotPassword(c echo.Context) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.authService.ForgotPassword(c.Request().Context(), req.Email)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Password reset email sent successfully", nil)
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset the user's password using the token from the password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body struct { Token string `json:"token" validate:"required"`; NewPassword string `json:"new_password" validate:"required,min=8"` } true "New password"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/reset-password [post]
func (h *authHandler) ResetPassword(c echo.Context) error {
	var req struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.authService.ResetPassword(c.Request().Context(), req.Token, req.NewPassword)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Password reset successfully", nil)
}

// VerifyEmail handles email verification
// @Summary Verify email address
// @Description Verify the user's email address using the token sent to their email
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/verify-email [get]
func (h *authHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Verification token is required")
	}

	err := h.authService.VerifyEmail(c.Request().Context(), token)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Email verified successfully", nil)
}

// ResendVerification handles resending email verification
// @Summary Resend email verification
// @Description Resend the email verification link to the user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body struct { Email string `json:"email" validate:"required,email"` } true "Email address"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/resend-verification [post]
func (h *authHandler) ResendVerification(c echo.Context) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.authService.ResendVerification(c.Request().Context(), req.Email)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Verification email sent successfully", nil)
}

// Helper function to get user ID from context
func getUserID(c echo.Context) (uint, error) {
	userIDStr := c.Get("user_id")
	if userIDStr == nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}
	
	userID, ok := userIDStr.(uint)
	if !ok {
		// Try to convert from string if it's stored as string
		if str, ok := userIDStr.(string); ok {
			id, err := strconv.ParseUint(str, 10, 32)
			if err != nil {
				return 0, echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
			}
			return uint(id), nil
		}
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
	}
	
	return userID, nil
}
