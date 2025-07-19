package service

import (
	"context"
	"errors"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"gorm.io/gorm"
)

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, req *models.NotificationCreateRequest) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Data:    req.Data,
		IsRead:  false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *notificationService) Create(ctx context.Context, req *models.NotificationCreateRequest) (*models.NotificationResponse, error) {
	notification := &models.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Data:    req.Data,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	resp := notification.ToResponse()
	return &resp, nil
}

func (s *notificationService) GetByUser(ctx context.Context, userID uint, page, limit int) (*models.NotificationListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	notifications, total, err := s.notificationRepo.GetByUser(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []models.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, notification.ToResponse())
	}

	return &models.NotificationListResponse{
		Notifications: responses,
		Total:         total,
		UnreadCount:   unreadCount,
	}, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, userID, notificationID uint) error {
	// Verify notification belongs to user
	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("notification not found")
		}
		return err
	}

	if notification.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.notificationRepo.MarkAsRead(ctx, userID, notificationID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uint) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID uint) (int, error) {
	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	return int(count), err
}

func (s *notificationService) GetUserNotifications(ctx context.Context, userID uint, limit, offset int) ([]*models.Notification, error) {
	page := (offset / limit) + 1
	notifications, _, err := s.notificationRepo.GetByUser(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	var result []*models.Notification
	for _, notification := range notifications {
		result = append(result, &notification)
	}

	return result, nil
}

func (s *notificationService) GetUnreadNotifications(ctx context.Context, userID uint) ([]*models.Notification, error) {
	// Since we don't have a specific GetUnreadByUserID method, we'll use GetByUser and filter
	notifications, _, err := s.notificationRepo.GetByUser(ctx, userID, 1, 100)
	if err != nil {
		return nil, err
	}

	var result []*models.Notification
	for _, notification := range notifications {
		if !notification.IsRead {
			result = append(result, &notification)
		}
	}

	return result, nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, userID, notificationID uint) error {
	// Verify notification belongs to user
	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("notification not found")
		}
		return err
	}

	if notification.UserID != userID {
		return errors.New("unauthorized")
	}

	// Since we don't have a specific Delete method, we'll use DeleteOld with 0 days
	return s.notificationRepo.DeleteOld(ctx, userID, 0)
}

func (s *notificationService) GetNotificationCount(ctx context.Context, userID uint) (int, error) {
	_, count, err := s.notificationRepo.GetByUser(ctx, userID, 1, 1)
	return int(count), err
}
