package service

import (
	"context"
	"errors"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"gorm.io/gorm"
)

type wishlistService struct {
	wishlistRepo repository.WishlistRepository
	productRepo  repository.ProductRepository
}

func NewWishlistService(wishlistRepo repository.WishlistRepository, productRepo repository.ProductRepository) WishlistService {
	return &wishlistService{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
	}
}

func (s *wishlistService) AddToWishlist(ctx context.Context, userID uint, req *models.WishlistAddRequest) (*models.WishlistResponse, error) {
	// Check if product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Check if already in wishlist
	exists, err := s.wishlistRepo.IsInWishlist(ctx, userID, req.ProductID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("product already in wishlist")
	}

	// Add to wishlist
	wishlistItem := &models.Wishlist{
		UserID:    userID,
		ProductID: req.ProductID,
	}

	if err := s.wishlistRepo.Add(ctx, wishlistItem); err != nil {
		return nil, err
	}

	// Prepare response
	wishlistItem.Product = *product
	resp := wishlistItem.ToResponse()
	return &resp, nil
}

func (s *wishlistService) GetWishlist(ctx context.Context, userID uint) (*models.WishlistItemsResponse, error) {
	wishlistItems, err := s.wishlistRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var items []models.WishlistResponse
	for _, item := range wishlistItems {
		items = append(items, item.ToResponse())
	}

	return &models.WishlistItemsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *wishlistService) RemoveFromWishlist(ctx context.Context, userID, productID uint) error {
	return s.wishlistRepo.Remove(ctx, userID, productID)
}

func (s *wishlistService) IsInWishlist(ctx context.Context, userID, productID uint) (bool, error) {
	return s.wishlistRepo.IsInWishlist(ctx, userID, productID)
}

func (s *wishlistService) GetUserWishlist(ctx context.Context, userID uint) ([]*models.WishlistResponse, error) {
	wishlistItems, err := s.wishlistRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []*models.WishlistResponse
	for _, item := range wishlistItems {
		response := &models.WishlistResponse{
			ID:        item.ID,
			UserID:    item.UserID,
			ProductID: item.ProductID,
			CreatedAt: item.CreatedAt,
		}

		// Get product details
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err == nil {
			productResp := product.ToResponse()
			response.Product = &productResp
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *wishlistService) IsProductInWishlist(ctx context.Context, userID, productID uint) (bool, error) {
	return s.wishlistRepo.IsInWishlist(ctx, userID, productID)
}

func (s *wishlistService) ClearWishlist(ctx context.Context, userID uint) error {
	// Since we don't have a direct clear method, we'll get all items and remove them
	wishlistItems, err := s.wishlistRepo.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	for _, item := range wishlistItems {
		if err := s.wishlistRepo.Remove(ctx, userID, item.ProductID); err != nil {
			return err
		}
	}

	return nil
}
