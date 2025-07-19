package repository

import (
	"context"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByUser(ctx context.Context, userID uint, page, limit int) ([]models.Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	MarkAsRead(ctx context.Context, userID, notificationID uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	GetByID(ctx context.Context, id uint) (*models.Notification, error)
	DeleteOld(ctx context.Context, userID uint, days int) error
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) GetByUser(ctx context.Context, userID uint, page, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64
	
	offset := (page - 1) * limit
	
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error
	
	return notifications, total, err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, userID, notificationID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *notificationRepository) GetByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.WithContext(ctx).First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) DeleteOld(ctx context.Context, userID uint, days int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND created_at < NOW() - INTERVAL ? DAY", userID, days).
		Delete(&models.Notification{}).Error
}
