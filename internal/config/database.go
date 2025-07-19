package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Note: Auto-migration is disabled since we use SQL migration files
	// If you need to enable auto-migration for development, uncomment the following lines:
	// if err := AutoMigrate(db); err != nil {
	//     return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	// }

	return db, nil
}

func InitRedis(cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return client, nil
}

func AutoMigrate(db *gorm.DB) error {
	// Import your models here
	return db.AutoMigrate(
		&models.User{},
		&models.PasswordResetToken{},
		&models.EmailVerificationToken{},
		&models.Category{},
		&models.Product{},
		&models.ProductImage{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
		&models.Review{},
		&models.ReviewHelpful{},
		&models.Wishlist{},
		&models.Notification{},
	)
}
