package repository

import (
	"context"
	"time"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetByStatus(ctx context.Context, status models.OrderStatus, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *orderRepository) UpdateTrackingNumber(ctx context.Context, id uint, trackingNumber string) error {
	return r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("tracking_number", trackingNumber).Error
}

func (r *orderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Order{}, id).Error
}

func (r *orderRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).Count(&count).Error
	return count, err
}

func (r *orderRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *orderRepository) CountByStatus(ctx context.Context, status models.OrderStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *orderRepository) GetTotalRevenue(ctx context.Context, startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("status = ?", models.OrderStatusDelivered).
		Select("COALESCE(SUM(total_amount), 0)")

	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Scan(&total).Error
	return total, err
}

func (r *orderRepository) GetOrdersBySellerID(ctx context.Context, sellerID uint, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.WithContext(ctx).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Where("products.seller_id = ?", sellerID).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Group("orders.id").
		Order("orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetRevenueBySellerID(ctx context.Context, sellerID uint, startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := r.db.WithContext(ctx).
		Model(&models.OrderItem{}).
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("products.seller_id = ? AND orders.status = ?", sellerID, models.OrderStatusDelivered).
		Select("COALESCE(SUM(order_items.price * order_items.quantity), 0)")

	if startDate != nil && endDate != nil {
		query = query.Where("orders.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Scan(&total).Error
	return total, err
}
