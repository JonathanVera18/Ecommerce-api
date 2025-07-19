package repository

import (
	"context"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type wishlistRepository struct {
	db *gorm.DB
}

type WishlistRepository interface {
	Add(ctx context.Context, wishlist *models.Wishlist) error
	GetByUser(ctx context.Context, userID uint) ([]models.Wishlist, error)
	Remove(ctx context.Context, userID, productID uint) error
	IsInWishlist(ctx context.Context, userID, productID uint) (bool, error)
	GetByUserAndProduct(ctx context.Context, userID, productID uint) (*models.Wishlist, error)
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) Add(ctx context.Context, wishlist *models.Wishlist) error {
	return r.db.WithContext(ctx).Create(wishlist).Error
}

func (r *wishlistRepository) GetByUser(ctx context.Context, userID uint) ([]models.Wishlist, error) {
	var wishlist []models.Wishlist
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.ProductImages").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&wishlist).Error
	return wishlist, err
}

func (r *wishlistRepository) Remove(ctx context.Context, userID, productID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Wishlist{}).Error
}

func (r *wishlistRepository) IsInWishlist(ctx context.Context, userID, productID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Wishlist{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}

func (r *wishlistRepository) GetByUserAndProduct(ctx context.Context, userID, productID uint) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&wishlist).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}
