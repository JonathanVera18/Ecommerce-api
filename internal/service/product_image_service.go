package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
)

type productImageService struct {
	productImageRepo repository.ProductImageRepository
	productRepo      repository.ProductRepository
}

func NewProductImageService(
	productImageRepo repository.ProductImageRepository,
	productRepo repository.ProductRepository,
) ProductImageService {
	return &productImageService{
		productImageRepo: productImageRepo,
		productRepo:      productRepo,
	}
}

func (s *productImageService) AddProductImage(ctx context.Context, productID uint, imageReq *models.ProductImageRequest) (*models.ProductImage, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Create product image
	productImage := &models.ProductImage{
		ProductID: productID,
		URL:       imageReq.URL,
		AltText:   imageReq.AltText,
		SortOrder: imageReq.SortOrder,
		IsPrimary: imageReq.IsPrimary,
	}

	// If this is set as primary, ensure no other image is primary
	if imageReq.IsPrimary {
		if err := s.productImageRepo.SetPrimary(ctx, productID, 0); err != nil {
			return nil, err
		}
	}

	if err := s.productImageRepo.Create(ctx, productImage); err != nil {
		return nil, err
	}

	return productImage, nil
}

func (s *productImageService) GetProductImages(ctx context.Context, productID uint) ([]models.ProductImage, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return s.productImageRepo.GetByProductID(ctx, productID)
}

func (s *productImageService) GetProductImage(ctx context.Context, imageID uint) (*models.ProductImage, error) {
	return s.productImageRepo.GetByID(ctx, imageID)
}

func (s *productImageService) UpdateProductImage(ctx context.Context, imageID uint, imageReq *models.ProductImageRequest) (*models.ProductImage, error) {
	// Get existing image
	existingImage, err := s.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingImage.URL = imageReq.URL
	existingImage.AltText = imageReq.AltText
	existingImage.SortOrder = imageReq.SortOrder

	// Handle primary image logic
	if imageReq.IsPrimary && !existingImage.IsPrimary {
		if err := s.productImageRepo.SetPrimary(ctx, existingImage.ProductID, imageID); err != nil {
			return nil, err
		}
		existingImage.IsPrimary = true
	} else if !imageReq.IsPrimary && existingImage.IsPrimary {
		existingImage.IsPrimary = false
	}

	if err := s.productImageRepo.Update(ctx, existingImage); err != nil {
		return nil, err
	}

	return existingImage, nil
}

func (s *productImageService) DeleteProductImage(ctx context.Context, imageID uint) error {
	// Get existing image to verify it exists
	_, err := s.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return err
	}

	return s.productImageRepo.Delete(ctx, imageID)
}

func (s *productImageService) SetPrimaryImage(ctx context.Context, productID uint, imageID uint) error {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Verify image exists and belongs to product
	image, err := s.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return errors.New("image not found")
	}

	if image.ProductID != productID {
		return errors.New("image does not belong to the specified product")
	}

	return s.productImageRepo.SetPrimary(ctx, productID, imageID)
}

func (s *productImageService) GetPrimaryImage(ctx context.Context, productID uint) (*models.ProductImage, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return s.productImageRepo.GetPrimaryImage(ctx, productID)
}

func (s *productImageService) UpdateImageOrder(ctx context.Context, productID uint, imageID uint, sortOrder int) error {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Verify image exists and belongs to product
	image, err := s.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return errors.New("image not found")
	}

	if image.ProductID != productID {
		return errors.New("image does not belong to the specified product")
	}

	return s.productImageRepo.UpdateSortOrder(ctx, productID, imageID, sortOrder)
}

func (s *productImageService) BulkAddImages(ctx context.Context, productID uint, imageReqs []models.ProductImageRequest) ([]models.ProductImage, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if len(imageReqs) == 0 {
		return []models.ProductImage{}, nil
	}

	// Convert requests to models
	var images []models.ProductImage
	hasPrimary := false
	
	for _, req := range imageReqs {
		image := models.ProductImage{
			ProductID: productID,
			URL:       req.URL,
			AltText:   req.AltText,
			SortOrder: req.SortOrder,
			IsPrimary: req.IsPrimary,
		}

		// Ensure only one primary image
		if req.IsPrimary {
			if hasPrimary {
				image.IsPrimary = false
			} else {
				hasPrimary = true
			}
		}

		images = append(images, image)
	}

	// If a primary image is being added, clear existing primary
	if hasPrimary {
		if err := s.productImageRepo.SetPrimary(ctx, productID, 0); err != nil {
			return nil, err
		}
	}

	if err := s.productImageRepo.BulkCreate(ctx, images); err != nil {
		return nil, err
	}

	return images, nil
}

func (s *productImageService) ReplaceProductImages(ctx context.Context, productID uint, imageReqs []models.ProductImageRequest) ([]models.ProductImage, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Delete existing images
	if err := s.productImageRepo.DeleteByProductID(ctx, productID); err != nil {
		return nil, fmt.Errorf("failed to delete existing images: %w", err)
	}

	// Add new images
	return s.BulkAddImages(ctx, productID, imageReqs)
}
