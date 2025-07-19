package main

import (
	"log"
	"os"

	"github.com/JonathanVera18/ecommerce-api/internal/config"
	"github.com/JonathanVera18/ecommerce-api/internal/handler"
	"github.com/JonathanVera18/ecommerce-api/internal/middleware"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	
	"github.com/JonathanVera18/ecommerce-api/pkg/payment"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// @title E-commerce API
// @version 1.0
// @description A comprehensive e-commerce backend API
// @contact.name API Support
// @contact.email support@ecommerce.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis
	redisClient, err := config.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize external services
	
	paymentService := payment.NewStripeService(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	wishlistRepo := repository.NewWishlistRepository(db)
	cartRepo := repository.NewCartRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	productImageRepo := repository.NewProductImageRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg, redisClient)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo, reviewRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, userRepo, paymentService)
	reviewService := service.NewReviewService(reviewRepo, productRepo, userRepo)
	categoryService := service.NewCategoryService(categoryRepo, productRepo)
	wishlistService := service.NewWishlistService(wishlistRepo, productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	productImageService := service.NewProductImageService(productImageRepo, productRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, authService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	adminHandler := handler.NewAdminHandler(userService, productService, orderService, reviewService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	wishlistHandler := handler.NewWishlistHandler(wishlistService)
	cartHandler := handler.NewCartHandler(cartService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	fileUploadHandler := handler.NewFileUploadHandler("uploads")
	productImageHandler := handler.NewProductImageHandler(productImageService)

	// Initialize Echo
	e := echo.New()

	// Global middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Logging())

	// Routes
	handler.SetupRoutes(e, &handler.Handlers{
		Auth:         authHandler,
		User:         userHandler,
		Product:      productHandler,
		Order:        orderHandler,
		Review:       reviewHandler,
		Admin:        adminHandler,
		Category:     categoryHandler,
		Wishlist:     wishlistHandler,
		Cart:         cartHandler,
		Notification: notificationHandler,
		FileUpload:   fileUploadHandler,
		ProductImage: productImageHandler,
	}, authService)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
