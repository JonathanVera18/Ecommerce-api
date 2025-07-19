package models

import (
	"time"
)

// Review represents a product review
type Review struct {
	BaseModel
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	User      User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrderID   *uint   `json:"order_id,omitempty"` // Optional: link to order for verified purchases
	Order     *Order  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	
	Rating  int    `json:"rating" gorm:"not null" validate:"required,min=1,max=5"`
	Title   string `json:"title" gorm:"type:varchar(255)" validate:"max=255"`
	Comment string `json:"comment" gorm:"type:text" validate:"required,min=10,max=2000"`
	
	// Review status
	IsVerified bool `json:"is_verified" gorm:"default:false"` // Verified purchase
	IsApproved bool `json:"is_approved" gorm:"default:true"`  // Moderation
	
	// Helpful votes
	HelpfulCount    int `json:"helpful_count" gorm:"default:0"`
	NotHelpfulCount int `json:"not_helpful_count" gorm:"default:0"`
	
	// Response from seller/admin
	SellerResponse   *string    `json:"seller_response,omitempty" gorm:"type:text"`
	SellerResponseAt *time.Time `json:"seller_response_at,omitempty"`
	ResponseBy       *uint      `json:"response_by,omitempty"`
	
	// Relationships
	ReviewHelpful []ReviewHelpful `json:"-" gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE"`
}

// ReviewHelpful represents helpful votes for reviews
type ReviewHelpful struct {
	BaseModel
	ReviewID uint `json:"review_id" gorm:"not null"`
	UserID   uint `json:"user_id" gorm:"not null"`
	IsHelpful bool `json:"is_helpful" gorm:"not null"` // true for helpful, false for not helpful
	
	// Unique constraint
	Review User `json:"-" gorm:"foreignKey:ReviewID"`
	User   User `json:"-" gorm:"foreignKey:UserID"`
}

// ReviewCreateRequest represents the request to create a review
type ReviewCreateRequest struct {
	ProductID uint   `json:"product_id" validate:"required"`
	OrderID   *uint  `json:"order_id,omitempty"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Title     string `json:"title" validate:"max=255"`
	Comment   string `json:"comment" validate:"required,min=10,max=2000"`
}

// ReviewUpdateRequest represents the request to update a review
type ReviewUpdateRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Title   *string `json:"title,omitempty" validate:"omitempty,max=255"`
	Comment *string `json:"comment,omitempty" validate:"omitempty,min=10,max=2000"`
}

// ReviewListRequest represents the request to list reviews with filters
type ReviewListRequest struct {
	Page       int     `query:"page" validate:"min=1"`
	Limit      int     `query:"limit" validate:"min=1,max=100"`
	ProductID  *uint   `query:"product_id"`
	UserID     *uint   `query:"user_id"`
	Rating     *int    `query:"rating" validate:"omitempty,min=1,max=5"`
	IsVerified *bool   `query:"is_verified"`
	IsApproved *bool   `query:"is_approved"`
	DateFrom   *time.Time `query:"date_from"`
	DateTo     *time.Time `query:"date_to"`
	SortBy     string  `query:"sort_by" validate:"omitempty,oneof=created_at rating helpful_count"`
	SortOrder  string  `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// ReviewResponse represents a review response
type ReviewResponse struct {
	ID               uint                 `json:"id"`
	ProductID        uint                 `json:"product_id"`
	Product          *ProductBasicInfo    `json:"product,omitempty"`
	UserID           uint                 `json:"user_id"`
	User             *UserBasicInfo       `json:"user,omitempty"`
	OrderID          *uint                `json:"order_id,omitempty"`
	Rating           int                  `json:"rating"`
	Title            string               `json:"title"`
	Comment          string               `json:"comment"`
	IsVerified       bool                 `json:"is_verified"`
	IsApproved       bool                 `json:"is_approved"`
	HelpfulCount     int                  `json:"helpful_count"`
	NotHelpfulCount  int                  `json:"not_helpful_count"`
	SellerResponse   *string              `json:"seller_response,omitempty"`
	SellerResponseAt *time.Time           `json:"seller_response_at,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

// ProductBasicInfo represents basic product information for reviews
type ProductBasicInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

// UserBasicInfo represents basic user information for reviews
type UserBasicInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar,omitempty"`
}

// ReviewStatsResponse represents review statistics
type ReviewStatsResponse struct {
	TotalReviews    int64   `json:"total_reviews"`
	AverageRating   float64 `json:"average_rating"`
	RatingDistribution map[int]int64 `json:"rating_distribution"`
	VerifiedReviews int64   `json:"verified_reviews"`
	PendingReviews  int64   `json:"pending_reviews"`
}

// ProductReviewStats represents review statistics for a product
type ProductReviewStats struct {
	ProductID          uint                `json:"product_id"`
	TotalReviews       int                 `json:"total_reviews"`
	AverageRating      float64             `json:"average_rating"`
	RatingDistribution map[int]int         `json:"rating_distribution"`
	VerifiedReviews    int                 `json:"verified_reviews"`
}

// ReviewHelpfulRequest represents the request to mark a review as helpful
type ReviewHelpfulRequest struct {
	IsHelpful bool `json:"is_helpful" validate:"required"`
}

// SellerResponseRequest represents the request to add a seller response
type SellerResponseRequest struct {
	Response string `json:"response" validate:"required,min=10,max=1000"`
}

// Request models
type CreateReviewRequest struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment" validate:"required,min=10,max=2000"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment *string `json:"comment,omitempty" validate:"omitempty,min=10,max=2000"`
}

// Response models
type ReviewStats struct {
	AverageRating      float64        `json:"average_rating"`
	TotalReviews       int64          `json:"total_reviews"`
	RatingDistribution map[int]int64  `json:"rating_distribution"`
}

// ToResponse converts Review to ReviewResponse
func (r *Review) ToResponse() ReviewResponse {
	response := ReviewResponse{
		ID:               r.ID,
		ProductID:        r.ProductID,
		UserID:           r.UserID,
		OrderID:          r.OrderID,
		Rating:           r.Rating,
		Title:            r.Title,
		Comment:          r.Comment,
		IsVerified:       r.IsVerified,
		IsApproved:       r.IsApproved,
		HelpfulCount:     r.HelpfulCount,
		NotHelpfulCount:  r.NotHelpfulCount,
		SellerResponse:   r.SellerResponse,
		SellerResponseAt: r.SellerResponseAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
	
	// Add product info if loaded
	if r.Product.ID != 0 {
		response.Product = &ProductBasicInfo{
			ID:    r.Product.ID,
			Name:  r.Product.Name,
			Image: r.Product.GetPrimaryImage(),
		}
	}
	
	// Add user info if loaded
	if r.User.ID != 0 {
		avatar := ""
		if r.User.Avatar != nil {
			avatar = *r.User.Avatar
		}
		response.User = &UserBasicInfo{
			ID:        r.User.ID,
			FirstName: r.User.FirstName,
			LastName:  r.User.LastName,
			Avatar:    avatar,
		}
	}
	
	return response
}

// CanEdit checks if a review can be edited
func (r *Review) CanEdit(userID uint) bool {
	// Only the author can edit within 24 hours of creation
	if r.UserID != userID {
		return false
	}
	
	// Allow editing within 24 hours
	return time.Since(r.CreatedAt) < 24*time.Hour
}

// CanDelete checks if a review can be deleted
func (r *Review) CanDelete(userID uint, isAdmin bool) bool {
	// Admin can delete any review
	if isAdmin {
		return true
	}
	
	// Author can delete their own review
	return r.UserID == userID
}

// CanAddSellerResponse checks if a seller response can be added
func (r *Review) CanAddSellerResponse(userID uint, productSellerID uint, isAdmin bool) bool {
	// Admin can respond to any review
	if isAdmin {
		return true
	}
	
	// Product seller can respond to reviews of their products
	return userID == productSellerID && r.SellerResponse == nil
}

// GetHelpfulPercentage calculates the helpful percentage
func (r *Review) GetHelpfulPercentage() float64 {
	total := r.HelpfulCount + r.NotHelpfulCount
	if total == 0 {
		return 0
	}
	return (float64(r.HelpfulCount) / float64(total)) * 100
}

// IsRecentReview checks if the review is recent (within 30 days)
func (r *Review) IsRecentReview() bool {
	return time.Since(r.CreatedAt) <= 30*24*time.Hour
}
