package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id uint) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uint) error
	GetChildren(ctx context.Context, parentID uint) ([]models.Category, error)
	GetRootCategories(ctx context.Context) ([]models.Category, error)
	GetWithProductCount(ctx context.Context) ([]models.Category, error)
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	category.Slug = r.generateSlug(category.Name)
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Order("sort_order, name").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("slug = ?", slug).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) GetChildren(ctx context.Context, parentID uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sort_order, name").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetRootCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Preload("Children").
		Order("sort_order, name").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetWithProductCount(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Select("categories.*, COUNT(products.id) as product_count").
		Joins("LEFT JOIN products ON products.category_id = categories.id").
		Group("categories.id").
		Order("sort_order, name").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	
	// Check if slug exists and make it unique
	var count int64
	originalSlug := slug
	for i := 1; ; i++ {
		r.db.Model(&models.Category{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, i)
	}
	
	return slug
}
