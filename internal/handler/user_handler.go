package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type userHandler struct {
	userService service.UserService
	authService service.AuthService
}

// UserHandler type alias for the concrete user handler
type UserHandler = userHandler

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, authService service.AuthService) *UserHandler {
	return &userHandler{
		userService: userService,
		authService: authService,
	}
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/profile [get]
func (h *userHandler) GetProfile(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	user, err := h.userService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return utils.NotFoundError(c, "User not found")
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", user)
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UserUpdateRequest true "Profile update request"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /users/profile [put]
func (h *userHandler) UpdateProfile(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req models.UserUpdateRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.userService.UpdateProfile(c.Request().Context(), userID, &req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to update profile")
	}

	return utils.SuccessResponse(c, "Profile updated successfully", user)
}

// GetUsers handles listing users (admin only)
// @Summary List users
// @Description Get a list of users with pagination and role filtering
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param role query string false "Filter by role" Enums(customer, seller, admin)
// @Success 200 {object} models.Response{data=[]models.UserResponse,meta=models.PaginationMeta}
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /users [get]
func (h *userHandler) GetUsers(c echo.Context) error {
	page, limit := utils.PaginationParams(c)
	
	var role *models.UserRole
	if roleStr := c.QueryParam("role"); roleStr != "" {
		r := models.UserRole(roleStr)
		role = &r
	}

	users, total, err := h.userService.GetUsers(c.Request().Context(), page, limit, role)
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve users")
	}

	meta := utils.BuildPaginationMeta(page, limit, total)
	return utils.SuccessResponseWithMeta(c, "Users retrieved successfully", users, meta)
}

// GetUser handles getting user by ID (admin only)
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
func (h *userHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.BadRequestError(c, "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), uint(id))
	if err != nil {
		return utils.NotFoundError(c, "User not found")
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
}

// CreateUser handles creating a new user (admin only)
// @Summary Create user
// @Description Create a new user account
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "User creation request"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /users [post]
func (h *userHandler) CreateUser(c echo.Context) error {
	var req models.UserCreateRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			return utils.ConflictError(c, err.Error())
		}
		return utils.InternalServerError(c, "Failed to create user")
	}

	return utils.CreatedResponse(c, "User created successfully", user)
}

// UpdateUser handles updating a user (admin only)
// @Summary Update user
// @Description Update a specific user by their ID
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UserUpdateRequest true "User update request"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [put]
func (h *userHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.BadRequestError(c, "Invalid user ID")
	}

	var req models.UserUpdateRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), uint(id), &req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to update user")
	}

	return utils.SuccessResponse(c, "User updated successfully", user)
}

// DeleteUser handles deleting a user (admin only)
// @Summary Delete user
// @Description Delete a specific user by their ID
// @Tags admin
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func (h *userHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.BadRequestError(c, "Invalid user ID")
	}

	if err := h.userService.DeleteUser(c.Request().Context(), uint(id)); err != nil {
		return utils.InternalServerError(c, "Failed to delete user")
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}
