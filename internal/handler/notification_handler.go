package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// CreateNotification creates a new notification (admin only)
func (h *NotificationHandler) CreateNotification(c echo.Context) error {
	var req models.NotificationCreateRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	notification, err := h.notificationService.CreateNotification(c.Request().Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, "Notification created successfully", notification)
}

// GetUserNotifications retrieves user's notifications with pagination
func (h *NotificationHandler) GetUserNotifications(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	notifications, err := h.notificationService.GetUserNotifications(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Notifications retrieved successfully", notifications)
}

// GetUnreadNotifications retrieves user's unread notifications
func (h *NotificationHandler) GetUnreadNotifications(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	notifications, err := h.notificationService.GetUnreadNotifications(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Unread notifications retrieved successfully", notifications)
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
	}

	err = h.notificationService.MarkAsRead(c.Request().Context(), userID, uint(notificationID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Notification marked as read", nil)
}

// MarkAllAsRead marks all notifications as read for a user
func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	err := h.notificationService.MarkAllAsRead(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "All notifications marked as read", nil)
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
	}

	err = h.notificationService.DeleteNotification(c.Request().Context(), userID, uint(notificationID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Notification deleted successfully", nil)
}

// GetNotificationCount retrieves total notification count for a user
func (h *NotificationHandler) GetNotificationCount(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	count, err := h.notificationService.GetNotificationCount(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Notification count retrieved successfully", map[string]int{"count": count})
}

// GetUnreadCount retrieves unread notification count for a user
func (h *NotificationHandler) GetUnreadCount(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	count, err := h.notificationService.GetUnreadCount(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Unread notification count retrieved successfully", map[string]int{"count": count})
}
