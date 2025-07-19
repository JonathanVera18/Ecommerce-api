package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JonathanVera18/ecommerce-api/internal/config"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

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

	// Get admin credentials from environment variables
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// Validate required environment variables
	if adminEmail == "" {
		log.Fatal("ADMIN_EMAIL environment variable is required")
	}
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is required")
	}

	// Validate password strength
	if len(adminPassword) < 12 {
		log.Fatal("Admin password must be at least 12 characters long")
	}

	// Check if admin user already exists
	var existingUser models.User
	if err := db.Where("email = ?", adminEmail).First(&existingUser).Error; err == nil {
		fmt.Println("Admin user already exists with email:", adminEmail)
		return
	}

	// Create admin user
	admin := &models.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     adminEmail,
		Role:      models.RoleAdmin,
		IsActive:  true,
		IsVerified: true,
	}

	if err := admin.HashPassword(adminPassword); err != nil {
		log.Fatal("Failed to hash admin password:", err)
	}

	// Save admin user
	if err := db.Create(admin).Error; err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	fmt.Printf("Admin user created successfully!\n")
	fmt.Printf("Email: %s\n", adminEmail)
	fmt.Printf("Role: %s\n", admin.Role)
	fmt.Println("Admin password has been set securely.")
	fmt.Println("⚠️  IMPORTANT: Change the password after first login for additional security!")
}
}
