package models

import (
	"time"
)

// NotificationType represents notification types
type NotificationType string

const (
	NotificationTypeOrderCreated   NotificationType = "order_created"
	NotificationTypeOrderUpdated   NotificationType = "order_updated"
	NotificationTypeOrderShipped   NotificationType = "order_shipped"
	NotificationTypeOrderDelivered NotificationType = "order_delivered"
	NotificationTypeProductLowStock NotificationType = "product_low_stock"
	NotificationTypeReviewReceived NotificationType = "review_received"
	NotificationTypePasswordReset  NotificationType = "password_reset"
	NotificationTypeEmailVerified  NotificationType = "email_verified"
	NotificationTypeGeneral        NotificationType = "general"
)

// Notification represents a user notification
type Notification struct {
	BaseModel
	UserID    uint             `json:"user_id" gorm:"not null;index"`
	Type      NotificationType `json:"type" gorm:"type:varchar(50);not null"`
	Title     string           `json:"title" gorm:"type:varchar(255);not null"`
	Message   string           `json:"message" gorm:"type:text;not null"`
	Data      *string          `json:"data,omitempty" gorm:"type:json"` // JSON data for additional context
	IsRead    bool             `json:"is_read" gorm:"default:false"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	
	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// NotificationCreateRequest represents the request to create a notification
type NotificationCreateRequest struct {
	UserID  uint             `json:"user_id" validate:"required"`
	Type    NotificationType `json:"type" validate:"required"`
	Title   string           `json:"title" validate:"required,min=1,max=255"`
	Message string           `json:"message" validate:"required,min=1"`
	Data    *string          `json:"data,omitempty"`
}

// NotificationResponse represents the notification response
type NotificationResponse struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Data      *string          `json:"data,omitempty"`
	IsRead    bool             `json:"is_read"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// ToResponse converts Notification to NotificationResponse
func (n *Notification) ToResponse() NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Data:      n.Data,
		IsRead:    n.IsRead,
		ReadAt:    n.ReadAt,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// NotificationListResponse represents the notification list response
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	UnreadCount   int64                  `json:"unread_count"`
}
