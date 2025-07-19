package service

import (
	"context"
	"errors"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"gorm.io/gorm"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, productRepo repository.ProductRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

func (s *categoryService) Create(ctx context.Context, req *models.CategoryCreateRequest) (*models.CategoryResponse, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	// Reload with relations
	category, err := s.categoryRepo.GetByID(ctx, category.ID)
	if err != nil {
		return nil, err
	}

	resp := category.ToResponse()
	return &resp, nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]models.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []models.CategoryResponse
	for _, category := range categories {
		responses = append(responses, category.ToResponse())
	}
	return responses, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uint) (*models.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	resp := category.ToResponse()
	return &resp, nil
}

func (s *categoryService) Update(ctx context.Context, id uint, req *models.CategoryUpdateRequest) (*models.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.ImageURL != nil {
		category.ImageURL = req.ImageURL
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	// Reload with relations
	category, err = s.categoryRepo.GetByID(ctx, category.ID)
	if err != nil {
		return nil, err
	}

	resp := category.ToResponse()
	return &resp, nil
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Check if category has products
	// This would require a method in ProductRepository to count products by category
	// For now, we'll proceed with deletion

	// Check if category has children
	children, err := s.categoryRepo.GetChildren(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("cannot delete category with subcategories")
	}

	return s.categoryRepo.Delete(ctx, id)
}

func (s *categoryService) GetRootCategories(ctx context.Context) ([]models.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	var responses []models.CategoryResponse
	for _, category := range categories {
		responses = append(responses, category.ToResponse())
	}
	return responses, nil
}

func (s *categoryService) GetWithProductCount(ctx context.Context) ([]models.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []models.CategoryResponse
	for _, category := range categories {
		response := category.ToResponse()
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, req *models.CategoryCreateRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []*models.Category
	for _, category := range categories {
		result = append(result, &category)
	}

	return result, nil
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	return s.categoryRepo.GetBySlug(ctx, slug)
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uint, req *models.CategoryUpdateRequest) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.ImageURL != nil {
		category.ImageURL = req.ImageURL
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	return s.categoryRepo.Delete(ctx, id)
}

func (s *categoryService) GetCategoriesHierarchy(ctx context.Context) ([]*models.Category, error) {
	categories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	var result []*models.Category
	for _, category := range categories {
		result = append(result, &category)
	}

	return result, nil
}

func (s *categoryService) GetCategoryChildren(ctx context.Context, parentID uint) ([]*models.Category, error) {
	categories, err := s.categoryRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	var result []*models.Category
	for _, category := range categories {
		result = append(result, &category)
	}

	return result, nil
}
