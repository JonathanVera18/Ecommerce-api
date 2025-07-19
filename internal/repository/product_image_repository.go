package repository

import (
	"context"
	"errors"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type productImageRepository struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) ProductImageRepository {
	return &productImageRepository{db: db}
}

func (r *productImageRepository) Create(ctx context.Context, productImage *models.ProductImage) error {
	if err := r.db.WithContext(ctx).Create(productImage).Error; err != nil {
		return err
	}
	return nil
}

func (r *productImageRepository) GetByProductID(ctx context.Context, productID uint) ([]models.ProductImage, error) {
	var images []models.ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC, created_at ASC").
		Find(&images).Error
	return images, err
}

func (r *productImageRepository) GetByID(ctx context.Context, id uint) (*models.ProductImage, error) {
	var image models.ProductImage
	err := r.db.WithContext(ctx).First(&image, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product image not found")
		}
		return nil, err
	}
	return &image, nil
}

func (r *productImageRepository) Update(ctx context.Context, productImage *models.ProductImage) error {
	return r.db.WithContext(ctx).Save(productImage).Error
}

func (r *productImageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ProductImage{}, id).Error
}

func (r *productImageRepository) DeleteByProductID(ctx context.Context, productID uint) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&models.ProductImage{}).Error
}

func (r *productImageRepository) SetPrimary(ctx context.Context, productID uint, imageID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, unset all primary images for this product
		if err := tx.Model(&models.ProductImage{}).
			Where("product_id = ?", productID).
			Update("is_primary", false).Error; err != nil {
			return err
		}

		// Then set the specified image as primary
		if err := tx.Model(&models.ProductImage{}).
			Where("id = ? AND product_id = ?", imageID, productID).
			Update("is_primary", true).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *productImageRepository) GetPrimaryImage(ctx context.Context, productID uint) (*models.ProductImage, error) {
	var image models.ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_primary = ?", productID, true).
		First(&image).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("primary image not found")
		}
		return nil, err
	}
	return &image, nil
}

func (r *productImageRepository) UpdateSortOrder(ctx context.Context, productID uint, imageID uint, sortOrder int) error {
	return r.db.WithContext(ctx).Model(&models.ProductImage{}).
		Where("id = ? AND product_id = ?", imageID, productID).
		Update("sort_order", sortOrder).Error
}

func (r *productImageRepository) BulkCreate(ctx context.Context, productImages []models.ProductImage) error {
	if len(productImages) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(productImages, 100).Error
}
