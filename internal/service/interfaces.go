package service

import (
	"context"
	"time"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, token string) (string, error)
	Logout(ctx context.Context, userID uint) error
	GetCurrentUser(ctx context.Context, userID uint) (*models.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *models.PasswordChangeRequest) error
	ValidateToken(token string) (uint, error)
	GetJWTService() *utils.JWTService
	// New methods for password reset and email verification
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
	VerifyEmail(ctx context.Context, token string) error
	ResendVerification(ctx context.Context, email string) error
}

// UserService defines the interface for user operations
type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*models.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	GetUsers(ctx context.Context, page, limit int, role *models.UserRole) ([]models.UserResponse, int64, error)
	GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error)
	CreateUser(ctx context.Context, req *models.UserCreateRequest) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
	GetUserStats(ctx context.Context) (*models.UserStatsResponse, error)
}

// ProductService defines the interface for product operations
type ProductService interface {
	CreateProduct(ctx context.Context, req *models.CreateProductRequest, sellerID uint) (*models.Product, error)
	GetProduct(ctx context.Context, id uint) (*models.Product, error)
	GetProducts(ctx context.Context, req *models.GetProductsRequest) (*models.ProductListResponse, error)
	UpdateProduct(ctx context.Context, id uint, req *models.UpdateProductRequest, sellerID uint) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uint, sellerID uint) error
	UpdateStock(ctx context.Context, id uint, stock int, sellerID uint) error
	GetLowStockProducts(ctx context.Context, threshold int, sellerID *uint) ([]*models.Product, error)
	GetTopRatedProducts(ctx context.Context, limit int) ([]*models.Product, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]*models.Product, error)
	GetProductsByCategory(ctx context.Context, category string, limit, offset int) ([]*models.Product, error)
	UpdateProductRating(ctx context.Context, productID uint) error
}

// OrderService defines the interface for order operations
type OrderService interface {
	CreateOrder(ctx context.Context, req *models.CreateOrderRequest, userID uint) (*models.Order, error)
	GetOrder(ctx context.Context, id uint, userID uint, userRole models.UserRole) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID uint, limit, offset int) ([]*models.Order, error)
	GetAllOrders(ctx context.Context, limit, offset int) ([]*models.Order, error)
	GetOrdersByStatus(ctx context.Context, status models.OrderStatus, limit, offset int) ([]*models.Order, error)
	GetSellerOrders(ctx context.Context, sellerID uint, limit, offset int) ([]*models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status models.OrderStatus, userID uint, userRole models.UserRole) error
	ProcessPayment(ctx context.Context, orderID uint, paymentReq *models.PaymentRequest) (*models.PaymentResponse, error)
	CancelOrder(ctx context.Context, id uint, userID uint, userRole models.UserRole) error
	GetOrderAnalytics(ctx context.Context, sellerID *uint, startDate, endDate *time.Time) (*models.OrderAnalytics, error)
}

// ReviewService defines the interface for review operations
type ReviewService interface {
	CreateReview(ctx context.Context, req *models.CreateReviewRequest, userID uint) (*models.Review, error)
	GetReview(ctx context.Context, id uint) (*models.Review, error)
	GetProductReviews(ctx context.Context, productID uint, limit, offset int) ([]*models.Review, error)
	GetUserReviews(ctx context.Context, userID uint, limit, offset int) ([]*models.Review, error)
	UpdateReview(ctx context.Context, id uint, req *models.UpdateReviewRequest, userID uint) (*models.Review, error)
	DeleteReview(ctx context.Context, id uint, userID uint, userRole models.UserRole) error
	GetReviewsByRating(ctx context.Context, rating int, limit, offset int) ([]*models.Review, error)
	GetTopReviews(ctx context.Context, limit int) ([]*models.Review, error)
	GetRecentReviews(ctx context.Context, limit int) ([]*models.Review, error)
	GetProductReviewStats(ctx context.Context, productID uint) (*models.ReviewStats, error)
	CanUserReview(ctx context.Context, userID, productID uint) (bool, error)
}

// EmailService defines the interface for email operations
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, user *models.User) error
	SendOrderConfirmationEmail(ctx context.Context, user *models.User, order *models.Order) error
	SendOrderStatusUpdateEmail(ctx context.Context, user *models.User, order *models.Order) error
	SendPasswordResetEmail(ctx context.Context, user *models.User, resetToken string) error
	SendEmailVerificationEmail(ctx context.Context, user *models.User, verificationToken string) error
	SendLowStockAlert(ctx context.Context, seller *models.User, product *models.Product) error
	SendNewReviewNotification(ctx context.Context, seller *models.User, product *models.Product, review *models.Review) error
}

// CategoryService defines the interface for category operations
type CategoryService interface {
	CreateCategory(ctx context.Context, req *models.CategoryCreateRequest) (*models.Category, error)
	GetCategory(ctx context.Context, id uint) (*models.Category, error)
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error)
	UpdateCategory(ctx context.Context, id uint, req *models.CategoryUpdateRequest) (*models.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoriesHierarchy(ctx context.Context) ([]*models.Category, error)
	GetCategoryChildren(ctx context.Context, parentID uint) ([]*models.Category, error)
}

// WishlistService defines the interface for wishlist operations
type WishlistService interface {
	AddToWishlist(ctx context.Context, userID uint, req *models.WishlistAddRequest) (*models.WishlistResponse, error)
	RemoveFromWishlist(ctx context.Context, userID uint, productID uint) error
	GetUserWishlist(ctx context.Context, userID uint) ([]*models.WishlistResponse, error)
	IsProductInWishlist(ctx context.Context, userID uint, productID uint) (bool, error)
	ClearWishlist(ctx context.Context, userID uint) error
}

// CartService defines the interface for cart operations
type CartService interface {
	AddToCart(ctx context.Context, userID uint, req *models.CartAddRequest) (*models.CartResponse, error)
	UpdateCartItem(ctx context.Context, userID uint, productID uint, quantity int) (*models.CartResponse, error)
	RemoveFromCart(ctx context.Context, userID uint, productID uint) error
	GetUserCart(ctx context.Context, userID uint) ([]*models.CartResponse, error)
	GetCartTotal(ctx context.Context, userID uint) (float64, error)
	ClearCart(ctx context.Context, userID uint) error
	GetCartItemCount(ctx context.Context, userID uint) (int, error)
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	CreateNotification(ctx context.Context, req *models.NotificationCreateRequest) (*models.Notification, error)
	GetUserNotifications(ctx context.Context, userID uint, limit, offset int) ([]*models.Notification, error)
	GetUnreadNotifications(ctx context.Context, userID uint) ([]*models.Notification, error)
	MarkAsRead(ctx context.Context, userID uint, notificationID uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	DeleteNotification(ctx context.Context, userID uint, notificationID uint) error
	GetNotificationCount(ctx context.Context, userID uint) (int, error)
	GetUnreadCount(ctx context.Context, userID uint) (int, error)
}

// ProductImageService defines the interface for product image operations
type ProductImageService interface {
	AddProductImage(ctx context.Context, productID uint, imageReq *models.ProductImageRequest) (*models.ProductImage, error)
	GetProductImages(ctx context.Context, productID uint) ([]models.ProductImage, error)
	GetProductImage(ctx context.Context, imageID uint) (*models.ProductImage, error)
	UpdateProductImage(ctx context.Context, imageID uint, imageReq *models.ProductImageRequest) (*models.ProductImage, error)
	DeleteProductImage(ctx context.Context, imageID uint) error
	SetPrimaryImage(ctx context.Context, productID uint, imageID uint) error
	GetPrimaryImage(ctx context.Context, productID uint) (*models.ProductImage, error)
	UpdateImageOrder(ctx context.Context, productID uint, imageID uint, sortOrder int) error
	BulkAddImages(ctx context.Context, productID uint, imageReqs []models.ProductImageRequest) ([]models.ProductImage, error)
	ReplaceProductImages(ctx context.Context, productID uint, imageReqs []models.ProductImageRequest) ([]models.ProductImage, error)
}
