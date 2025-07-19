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

	// Create admin user
	adminEmail := "admin@example.com"
	adminPassword := "admin123"

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

	// Set password from environment or use default
	if envPassword := os.Getenv("ADMIN_PASSWORD"); envPassword != "" {
		adminPassword = envPassword
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
	fmt.Printf("Password: %s\n", adminPassword)
	fmt.Printf("Role: %s\n", admin.Role)
	fmt.Println("Please change the password after first login!")
}
