package repository

import (
	"context"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) GetByID(ctx context.Context, id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByProductID(ctx context.Context, productID uint, limit, offset int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Product").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetByRating(ctx context.Context, rating int, limit, offset int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.WithContext(ctx).
		Where("rating = ?", rating).
		Preload("User").
		Preload("Product").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *reviewRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, id).Error
}

func (r *reviewRepository) GetByUserAndProduct(ctx context.Context, userID, productID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Preload("User").
		Preload("Product").
		First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Review{}).Count(&count).Error
	return count, err
}

func (r *reviewRepository) CountByProductID(ctx context.Context, productID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count, err
}

func (r *reviewRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *reviewRepository) GetAverageRatingByProductID(ctx context.Context, productID uint) (float64, error) {
	var avgRating float64
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error
	return avgRating, err
}

func (r *reviewRepository) GetRatingDistribution(ctx context.Context, productID uint) (map[int]int64, error) {
	type RatingCount struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	}

	var results []RatingCount
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Order("rating").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	distribution := make(map[int]int64)
	for _, result := range results {
		distribution[result.Rating] = result.Count
	}

	return distribution, nil
}

func (r *reviewRepository) GetTopReviews(ctx context.Context, limit int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.WithContext(ctx).
		Where("rating >= ?", 4).
		Preload("User").
		Preload("Product").
		Order("rating DESC, created_at DESC").
		Limit(limit).
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetRecentReviews(ctx context.Context, limit int) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Order("created_at DESC").
		Limit(limit).
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) CheckUserCanReview(ctx context.Context, userID, productID uint) (bool, error) {
	// Check if user has purchased this product and order is delivered
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status = ?",
			userID, productID, models.OrderStatusDelivered).
		Count(&count).Error

	return count > 0, err
}
