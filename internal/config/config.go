package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	Database DatabaseConfig

	// Redis
	Redis RedisConfig

	// JWT
	JWT JWTConfig

	// Server
	Server ServerConfig

	// Email
	Email EmailConfig

	// Stripe
	Stripe StripeConfig

	// Application
	App AppConfig

	// File Upload
	Upload UploadConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type ServerConfig struct {
	Host string
	Port int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
}

type AppConfig struct {
	Environment string
	URL         string
	FrontendURL string
}

type UploadConfig struct {
	MaxFileSize int64
	UploadDir   string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional
	}

	config := &Config{}

	// Database configuration
	config.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		Name:     getEnv("DB_NAME", "ecommerce_db"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	// Redis configuration
	config.Redis = RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
	}

	// JWT configuration
	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY format: %w", err)
	}

	config.JWT = JWTConfig{
		Secret: getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		Expiry: jwtExpiry,
	}

	// Server configuration
	config.Server = ServerConfig{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: getEnvAsInt("SERVER_PORT", 8080),
	}

	// Email configuration
	config.Email = EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@ecommerce.com"),
	}

	// Stripe configuration
	config.Stripe = StripeConfig{
		SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
	}

	// Application configuration
	config.App = AppConfig{
		Environment: getEnv("APP_ENV", "development"),
		URL:         getEnv("APP_URL", "http://localhost:8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	// Upload configuration
	config.Upload = UploadConfig{
		MaxFileSize: getEnvAsInt64("MAX_FILE_SIZE", 10485760), // 10MB
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
