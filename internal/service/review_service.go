package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
)

type reviewService struct {
	reviewRepo  repository.ReviewRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
) ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, req *models.CreateReviewRequest, userID uint) (*models.Review, error) {
	// Validate user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Validate product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Check if user can review this product (has purchased it)
	canReview, err := s.reviewRepo.CheckUserCanReview(ctx, userID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to check review eligibility: %w", err)
	}

	if !canReview {
		return nil, errors.New("you can only review products you have purchased and received")
	}

	// Check if user has already reviewed this product
	existingReview, err := s.reviewRepo.GetByUserAndProduct(ctx, userID, req.ProductID)
	if err == nil && existingReview != nil {
		return nil, errors.New("you have already reviewed this product")
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	review := &models.Review{
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		User:      *user,
		Product:   *product,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Update product rating after creating review
	if err := s.updateProductRating(ctx, req.ProductID); err != nil {
		// Log error but don't fail the review creation
		fmt.Printf("Warning: failed to update product rating: %v\n", err)
	}

	return review, nil
}

func (s *reviewService) GetReview(ctx context.Context, id uint) (*models.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return review, nil
}

func (s *reviewService) GetProductReviews(ctx context.Context, productID uint, limit, offset int) ([]*models.Review, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	reviews, err := s.reviewRepo.GetByProductID(ctx, productID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get product reviews: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) GetUserReviews(ctx context.Context, userID uint, limit, offset int) ([]*models.Review, error) {
	// Validate user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	reviews, err := s.reviewRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) UpdateReview(ctx context.Context, id uint, req *models.UpdateReviewRequest, userID uint) (*models.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	// Check if user owns this review
	if review.UserID != userID {
		return nil, errors.New("unauthorized to update this review")
	}

	// Update fields if provided
	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return nil, errors.New("rating must be between 1 and 5")
		}
		review.Rating = *req.Rating
	}

	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	// Update product rating after updating review
	if err := s.updateProductRating(ctx, review.ProductID); err != nil {
		// Log error but don't fail the review update
		fmt.Printf("Warning: failed to update product rating: %v\n", err)
	}

	return review, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id uint, userID uint, userRole models.UserRole) error {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get review: %w", err)
	}

	// Check authorization
	if userRole != models.RoleAdmin && review.UserID != userID {
		return errors.New("unauthorized to delete this review")
	}

	productID := review.ProductID

	if err := s.reviewRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	// Update product rating after deleting review
	if err := s.updateProductRating(ctx, productID); err != nil {
		// Log error but don't fail the review deletion
		fmt.Printf("Warning: failed to update product rating: %v\n", err)
	}

	return nil
}

func (s *reviewService) GetReviewsByRating(ctx context.Context, rating int, limit, offset int) ([]*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	reviews, err := s.reviewRepo.GetByRating(ctx, rating, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews by rating: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) GetTopReviews(ctx context.Context, limit int) ([]*models.Review, error) {
	reviews, err := s.reviewRepo.GetTopReviews(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top reviews: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) GetRecentReviews(ctx context.Context, limit int) ([]*models.Review, error) {
	reviews, err := s.reviewRepo.GetRecentReviews(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent reviews: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) GetProductReviewStats(ctx context.Context, productID uint) (*models.ReviewStats, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Get average rating
	avgRating, err := s.reviewRepo.GetAverageRatingByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}

	// Get total review count
	totalReviews, err := s.reviewRepo.CountByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get review count: %w", err)
	}

	// Get rating distribution
	distribution, err := s.reviewRepo.GetRatingDistribution(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}

	return &models.ReviewStats{
		AverageRating:     avgRating,
		TotalReviews:      totalReviews,
		RatingDistribution: distribution,
	}, nil
}

func (s *reviewService) CanUserReview(ctx context.Context, userID, productID uint) (bool, error) {
	canReview, err := s.reviewRepo.CheckUserCanReview(ctx, userID, productID)
	if err != nil {
		return false, fmt.Errorf("failed to check review eligibility: %w", err)
	}

	// Also check if user has already reviewed
	if canReview {
		existingReview, err := s.reviewRepo.GetByUserAndProduct(ctx, userID, productID)
		if err == nil && existingReview != nil {
			return false, nil // User has already reviewed
		}
	}

	return canReview, nil
}

func (s *reviewService) updateProductRating(ctx context.Context, productID uint) error {
	// Get average rating
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
