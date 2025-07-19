package repository

import (
	"context"
	"fmt"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Preload("Reviews").
		Preload("Reviews.User").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Preload("Reviews").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (r *productRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Preload("Reviews").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (r *productRepository) GetBySellerID(ctx context.Context, sellerID uint, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Where("seller_id = ?", sellerID).
		Preload("Reviews").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	searchPattern := fmt.Sprintf("%%%s%%", query)
	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern).
		Preload("Reviews").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint, stock int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", id).
		Update("stock", stock).Error
}

func (r *productRepository) GetLowStock(ctx context.Context, threshold int) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Where("stock <= ?", threshold).
		Find(&products).Error
	return products, err
}

func (r *productRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *productRepository) CountByCategory(ctx context.Context, category string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("category = ?", category).
		Count(&count).Error
	return count, err
}

func (r *productRepository) GetTopRated(ctx context.Context, limit int) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Preload("Reviews").
		Order("average_rating DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

func (r *productRepository) UpdateRating(ctx context.Context, productID uint, averageRating float64, reviewCount int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"average_rating": averageRating,
			"review_count":   reviewCount,
		}).Error
}
