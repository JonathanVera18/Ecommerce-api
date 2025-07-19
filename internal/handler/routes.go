package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/middleware"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	
)

// Handlers contains all the handlers
type Handlers struct {
	Auth         *AuthHandler
	User         *UserHandler
	Product      *ProductHandler
	Order        *OrderHandler
	Review       *ReviewHandler
	Admin        *AdminHandler
	Category     *CategoryHandler
	Wishlist     *WishlistHandler
	Cart         *CartHandler
	Notification *NotificationHandler
	FileUpload   *FileUploadHandler
	ProductImage *ProductImageHandler
}

// SetupRoutes configures all the application routes
func SetupRoutes(e *echo.Echo, handlers *Handlers, authService service.AuthService) {
	// Get JWT service from auth service
	jwtService := authService.GetJWTService()

	// API version group
	api := e.Group("/api/v1")

	// Auth routes (no authentication required)
	auth := api.Group("/auth")
	auth.POST("/register", handlers.Auth.Register)
	auth.POST("/login", handlers.Auth.Login)
	auth.POST("/refresh", handlers.Auth.RefreshToken)
	auth.POST("/logout", handlers.Auth.Logout, middleware.JWTAuth(jwtService))
	auth.GET("/profile", handlers.Auth.GetProfile, middleware.JWTAuth(jwtService))
	auth.POST("/change-password", handlers.Auth.ChangePassword, middleware.JWTAuth(jwtService))
	auth.POST("/forgot-password", handlers.Auth.ForgotPassword)
	auth.POST("/reset-password", handlers.Auth.ResetPassword)
	auth.GET("/verify-email", handlers.Auth.VerifyEmail)
	auth.POST("/resend-verification", handlers.Auth.ResendVerification)

	// User routes
	users := api.Group("/users")
	users.GET("/me", handlers.User.GetProfile, middleware.JWTAuth(jwtService))
	users.GET("/profile", handlers.User.GetProfile, middleware.JWTAuth(jwtService))
	users.PUT("/profile", handlers.User.UpdateProfile, middleware.JWTAuth(jwtService))
	users.GET("", handlers.User.GetUsers, middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))
	users.GET("/:id", handlers.User.GetUser, middleware.JWTAuth(jwtService))
	users.POST("", handlers.User.CreateUser, middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))
	users.PUT("/:id", handlers.User.UpdateUser, middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))
	users.DELETE("/:id", handlers.User.DeleteUser, middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))

	// Product routes
	products := api.Group("/products")
	products.GET("", handlers.Product.GetProducts)
	products.GET("/:id", handlers.Product.GetProduct)
	products.POST("", handlers.Product.CreateProduct, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.PUT("/:id", handlers.Product.UpdateProduct, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.DELETE("/:id", handlers.Product.DeleteProduct, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.PUT("/:id/stock", handlers.Product.UpdateStock, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.GET("/low-stock", handlers.Product.GetLowStockProducts, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.GET("/top-rated", handlers.Product.GetTopRatedProducts)
	products.GET("/search", handlers.Product.SearchProducts)
	products.GET("/category/:category", handlers.Product.GetProductsByCategory)

	// Product reviews
	products.GET("/:product_id/reviews", handlers.Review.GetProductReviews)
	products.GET("/:product_id/reviews/stats", handlers.Review.GetProductReviewStats)
	products.GET("/:product_id/can-review", handlers.Review.CanUserReview, middleware.JWTAuth(jwtService))

	// Product images
	products.GET("/:product_id/images", handlers.ProductImage.GetProductImages)
	products.POST("/:product_id/images", handlers.ProductImage.AddProductImage, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.GET("/:product_id/images/primary", handlers.ProductImage.GetPrimaryImage)
	products.GET("/:product_id/images/:image_id", handlers.ProductImage.GetProductImage)
	products.PUT("/:product_id/images/:image_id", handlers.ProductImage.UpdateProductImage, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.DELETE("/:product_id/images/:image_id", handlers.ProductImage.DeleteProductImage, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.PUT("/:product_id/images/:image_id/primary", handlers.ProductImage.SetPrimaryImage, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.PUT("/:product_id/images/:image_id/order", handlers.ProductImage.UpdateImageOrder, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.POST("/:product_id/images/bulk", handlers.ProductImage.BulkAddImages, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	products.PUT("/:product_id/images/replace", handlers.ProductImage.ReplaceProductImages, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))

	// Order routes
	orders := api.Group("/orders")
	orders.POST("", handlers.Order.CreateOrder, middleware.JWTAuth(jwtService))
	orders.GET("/my", handlers.Order.GetUserOrders, middleware.JWTAuth(jwtService))
	orders.GET("/:id", handlers.Order.GetOrder, middleware.JWTAuth(jwtService))
	orders.PUT("/:id/status", handlers.Order.UpdateOrderStatus, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	orders.POST("/:id/payment", handlers.Order.ProcessPayment, middleware.JWTAuth(jwtService))
	orders.PUT("/:id/cancel", handlers.Order.CancelOrder, middleware.JWTAuth(jwtService))
	orders.GET("/status/:status", handlers.Order.GetOrdersByStatus, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))
	orders.GET("/analytics", handlers.Order.GetOrderAnalytics, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))

	// Seller routes
	seller := api.Group("/seller")
	seller.GET("/orders", handlers.Order.GetSellerOrders, middleware.JWTAuth(jwtService), middleware.RequireRole("seller", "admin"))

	// Review routes
	reviews := api.Group("/reviews")
	reviews.POST("", handlers.Review.CreateReview, middleware.JWTAuth(jwtService))
	reviews.GET("/my", handlers.Review.GetUserReviews, middleware.JWTAuth(jwtService))
	reviews.GET("/:id", handlers.Review.GetReview)
	reviews.PUT("/:id", handlers.Review.UpdateReview, middleware.JWTAuth(jwtService))
	reviews.DELETE("/:id", handlers.Review.DeleteReview, middleware.JWTAuth(jwtService))
	reviews.GET("/rating/:rating", handlers.Review.GetReviewsByRating)
	reviews.GET("/top", handlers.Review.GetTopReviews)
	reviews.GET("/recent", handlers.Review.GetRecentReviews)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))
	admin.GET("/dashboard", handlers.Admin.GetDashboardStats)
	admin.GET("/orders", handlers.Order.GetAllOrders)
	admin.GET("/orders/:id", handlers.Admin.GetOrderDetails)
	admin.PUT("/users/:id", handlers.Admin.ManageUser)
	admin.GET("/health", handlers.Admin.GetSystemHealth)
	
	// Admin analytics
	adminAnalytics := admin.Group("/analytics")
	adminAnalytics.GET("/sales", handlers.Admin.GetSalesAnalytics)
	adminAnalytics.GET("/users", handlers.Admin.GetUserAnalytics)
	adminAnalytics.GET("/products", handlers.Admin.GetProductAnalytics)
	adminAnalytics.GET("/reviews", handlers.Admin.GetReviewAnalytics)

	// Category routes
	categories := api.Group("/categories")
	categories.GET("", handlers.Category.GetAllCategories)
	categories.GET("/:id", handlers.Category.GetCategory)
	categories.GET("/slug/:slug", handlers.Category.GetCategoryBySlug)
	categories.GET("/hierarchy", handlers.Category.GetCategoriesHierarchy)
	categories.GET("/:parentId/children", handlers.Category.GetCategoryChildren)
	categories.POST("", handlers.Category.CreateCategory, middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))
	categories.PUT("/:id", handlers.Category.UpdateCategory, middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))
	categories.DELETE("/:id", handlers.Category.DeleteCategory, middleware.JWTAuth(jwtService), middleware.RequireRole("admin"))

	// Wishlist routes
	wishlist := api.Group("/wishlist")
	wishlist.Use(middleware.JWTAuth(jwtService))
	wishlist.POST("", handlers.Wishlist.AddToWishlist)
	wishlist.GET("", handlers.Wishlist.GetUserWishlist)
	wishlist.DELETE("/:productId", handlers.Wishlist.RemoveFromWishlist)
	wishlist.GET("/:productId/check", handlers.Wishlist.IsProductInWishlist)
	wishlist.DELETE("", handlers.Wishlist.ClearWishlist)

	// Cart routes
	cart := api.Group("/cart")
	cart.Use(middleware.JWTAuth(jwtService))
	cart.POST("", handlers.Cart.AddToCart)
	cart.GET("", handlers.Cart.GetUserCart)
	cart.PUT("/:productId", handlers.Cart.UpdateCartItem)
	cart.DELETE("/:productId", handlers.Cart.RemoveFromCart)
	cart.GET("/total", handlers.Cart.GetCartTotal)
	cart.GET("/count", handlers.Cart.GetCartItemCount)
	cart.DELETE("", handlers.Cart.ClearCart)

	// Notification routes
	notifications := api.Group("/notifications")
	notifications.Use(middleware.JWTAuth(jwtService))
	notifications.GET("", handlers.Notification.GetUserNotifications)
	notifications.GET("/unread", handlers.Notification.GetUnreadNotifications)
	notifications.PUT("/:id/read", handlers.Notification.MarkAsRead)
	notifications.PUT("/read-all", handlers.Notification.MarkAllAsRead)
	notifications.DELETE("/:id", handlers.Notification.DeleteNotification)
	notifications.GET("/count", handlers.Notification.GetNotificationCount)
	notifications.GET("/unread-count", handlers.Notification.GetUnreadCount)
	notifications.POST("", handlers.Notification.CreateNotification, middleware.RequireRole("admin"))

	// File upload routes
	uploads := api.Group("/uploads")
	uploads.POST("", handlers.FileUpload.UploadFile, middleware.JWTAuth(jwtService))
	uploads.GET("/my-files", handlers.FileUpload.GetUserFiles, middleware.JWTAuth(jwtService))
	uploads.DELETE("/:filename", handlers.FileUpload.DeleteFile, middleware.JWTAuth(jwtService))
	uploads.GET("/user_:userId/:filename", handlers.FileUpload.ServeFile)
}
