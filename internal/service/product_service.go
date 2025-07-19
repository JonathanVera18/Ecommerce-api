package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
)

type productService struct {
	productRepo repository.ProductRepository
	reviewRepo  repository.ReviewRepository
}

func NewProductService(productRepo repository.ProductRepository, reviewRepo repository.ReviewRepository) ProductService {
	return &productService{
		productRepo: productRepo,
		reviewRepo:  reviewRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *models.CreateProductRequest, sellerID uint) (*models.Product, error) {
	if req.Price <= 0 {
		return nil, errors.New("product price must be greater than 0")
	}

	if req.Stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      req.Images,
		SellerID:    sellerID,
		IsActive:    true,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProducts(ctx context.Context, req *models.GetProductsRequest) (*models.ProductListResponse, error) {
	var products []*models.Product
	var err error
	var total int64

	switch {
	case req.Category != "":
		products, err = s.productRepo.GetByCategory(ctx, req.Category, req.Limit, req.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get products by category: %w", err)
		}
		total, err = s.productRepo.CountByCategory(ctx, req.Category)
	case req.SellerID != nil:
		products, err = s.productRepo.GetBySellerID(ctx, *req.SellerID, req.Limit, req.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get products by seller: %w", err)
		}
		// For seller products, we could implement a specific count method
		total, err = s.productRepo.Count(ctx)
	case req.Search != "":
		products, err = s.productRepo.Search(ctx, req.Search, req.Limit, req.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to search products: %w", err)
		}
		total, err = s.productRepo.Count(ctx) // Simplified for now
	default:
		products, err = s.productRepo.GetAll(ctx, req.Limit, req.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get all products: %w", err)
		}
		total, err = s.productRepo.Count(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get product count: %w", err)
	}

	return &models.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uint, req *models.UpdateProductRequest, sellerID uint) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized to update this product")
	}

	// Update fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, errors.New("product price must be greater than 0")
		}
		product.Price = *req.Price
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("product stock cannot be negative")
		}
		product.Stock = *req.Stock
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint, sellerID uint) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product.SellerID != sellerID {
		return errors.New("unauthorized to delete this product")
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *productService) UpdateStock(ctx context.Context, id uint, stock int, sellerID uint) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product.SellerID != sellerID {
		return errors.New("unauthorized to update this product's stock")
	}

	if stock < 0 {
		return errors.New("stock cannot be negative")
	}

	if err := s.productRepo.UpdateStock(ctx, id, stock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

func (s *productService) GetLowStockProducts(ctx context.Context, threshold int, sellerID *uint) ([]*models.Product, error) {
	products, err := s.productRepo.GetLowStock(ctx, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}

	// Filter by seller if specified
	if sellerID != nil {
		var filteredProducts []*models.Product
		for _, product := range products {
			if product.SellerID == *sellerID {
				filteredProducts = append(filteredProducts, product)
			}
		}
		return filteredProducts, nil
	}

	return products, nil
}

func (s *productService) GetTopRatedProducts(ctx context.Context, limit int) ([]*models.Product, error) {
	products, err := s.productRepo.GetTopRated(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top rated products: %w", err)
	}

	return products, nil
}

func (s *productService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*models.Product, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("search query cannot be empty")
	}

	products, err := s.productRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	return products, nil
}

func (s *productService) GetProductsByCategory(ctx context.Context, category string, limit, offset int) ([]*models.Product, error) {
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("category cannot be empty")
	}

	products, err := s.productRepo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	return products, nil
}

func (s *productService) UpdateProductRating(ctx context.Context, productID uint) error {
	// Get average rating from reviews
	avgRating, err := s.reviewRepo.GetAverageRatingByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get average rating: %w", err)
	}

	// Get review count
	reviewCount, err := s.reviewRepo.CountByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get review count: %w", err)
	}

	// Update product rating
	if err := s.productRepo.UpdateRating(ctx, productID, avgRating, int(reviewCount)); err != nil {
		return fmt.Errorf("failed to update product rating: %w", err)
	}

	return nil
}
