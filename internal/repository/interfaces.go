package repository

import (
	"context"
	"time"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, limit int, role *models.UserRole) ([]models.User, int64, error)
	UpdateLastLogin(ctx context.Context, id uint) error
	GetStats(ctx context.Context) (*models.UserStatsResponse, error)
	CreatePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenStr string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenStr string) error
	CreateEmailVerificationToken(ctx context.Context, token *models.EmailVerificationToken) error
	GetEmailVerificationToken(ctx context.Context, tokenStr string) (*models.EmailVerificationToken, error)
	MarkEmailVerificationTokenUsed(ctx context.Context, tokenStr string) error
	MarkEmailVerified(ctx context.Context, userID uint) error
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*models.Product, error)
	GetBySellerID(ctx context.Context, sellerID uint, limit, offset int) ([]*models.Product, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	UpdateStock(ctx context.Context, id uint, stock int) error
	GetLowStock(ctx context.Context, threshold int) ([]*models.Product, error)
	Count(ctx context.Context) (int64, error)
	CountByCategory(ctx context.Context, category string) (int64, error)
	GetTopRated(ctx context.Context, limit int) ([]*models.Product, error)
	UpdateRating(ctx context.Context, productID uint, averageRating float64, reviewCount int) error
}

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*models.Order, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Order, error)
	GetByStatus(ctx context.Context, status models.OrderStatus, limit, offset int) ([]*models.Order, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error
	UpdateTrackingNumber(ctx context.Context, id uint, trackingNumber string) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	CountByStatus(ctx context.Context, status models.OrderStatus) (int64, error)
	GetTotalRevenue(ctx context.Context, startDate, endDate *time.Time) (float64, error)
	GetOrdersBySellerID(ctx context.Context, sellerID uint, limit, offset int) ([]*models.Order, error)
	GetRevenueBySellerID(ctx context.Context, sellerID uint, startDate, endDate *time.Time) (float64, error)
}

// ReviewRepository defines the interface for review data operations
type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id uint) (*models.Review, error)
	GetByProductID(ctx context.Context, productID uint, limit, offset int) ([]*models.Review, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*models.Review, error)
	GetByRating(ctx context.Context, rating int, limit, offset int) ([]*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uint) error
	GetByUserAndProduct(ctx context.Context, userID, productID uint) (*models.Review, error)
	Count(ctx context.Context) (int64, error)
	CountByProductID(ctx context.Context, productID uint) (int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	GetAverageRatingByProductID(ctx context.Context, productID uint) (float64, error)
	GetRatingDistribution(ctx context.Context, productID uint) (map[int]int64, error)
	GetTopReviews(ctx context.Context, limit int) ([]*models.Review, error)
	GetRecentReviews(ctx context.Context, limit int) ([]*models.Review, error)
	CheckUserCanReview(ctx context.Context, userID, productID uint) (bool, error)
}

// ProductImageRepository defines the interface for product image data operations
type ProductImageRepository interface {
	Create(ctx context.Context, productImage *models.ProductImage) error
	GetByProductID(ctx context.Context, productID uint) ([]models.ProductImage, error)
	GetByID(ctx context.Context, id uint) (*models.ProductImage, error)
	Update(ctx context.Context, productImage *models.ProductImage) error
	Delete(ctx context.Context, id uint) error
	DeleteByProductID(ctx context.Context, productID uint) error
	SetPrimary(ctx context.Context, productID uint, imageID uint) error
	GetPrimaryImage(ctx context.Context, productID uint) (*models.ProductImage, error)
	UpdateSortOrder(ctx context.Context, productID uint, imageID uint, sortOrder int) error
	BulkCreate(ctx context.Context, productImages []models.ProductImage) error
}

// UserStatsResponse represents user statistics (defined here to avoid circular imports)
type UserStatsResponse struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	VerifiedUsers  int64 `json:"verified_users"`
	Customers      int64 `json:"customers"`
	Sellers        int64 `json:"sellers"`
	Admins         int64 `json:"admins"`
	NewUsersToday  int64 `json:"new_users_today"`
	NewUsersWeek   int64 `json:"new_users_week"`
	NewUsersMonth  int64 `json:"new_users_month"`
}
