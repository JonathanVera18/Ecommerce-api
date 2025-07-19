package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/JonathanVera18/ecommerce-api/internal/config"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
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

	// Seed data
	log.Println("Starting database seeding...")

	// Create sample users
	createSampleUsers(db)
	
	// Create sample products
	createSampleProducts(db)
	
	// Create sample orders
	createSampleOrders(db)
	
	// Create sample reviews
	createSampleReviews(db)

	log.Println("Database seeding completed successfully!")
}

func createSampleUsers(db *gorm.DB) {
	users := []*models.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Role:      models.RoleCustomer,
			IsActive:  true,
			IsVerified: true,
		},
		{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			Role:      models.RoleCustomer,
			IsActive:  true,
			IsVerified: true,
		},
		{
			FirstName: "Bob",
			LastName:  "Wilson",
			Email:     "bob.wilson@example.com",
			Role:      models.RoleSeller,
			IsActive:  true,
			IsVerified: true,
			StoreName: "Bob's Electronics",
			StoreDescription: "Quality electronics at great prices",
		},
		{
			FirstName: "Alice",
			LastName:  "Brown",
			Email:     "alice.brown@example.com",
			Role:      models.RoleSeller,
			IsActive:  true,
			IsVerified: true,
			StoreName: "Alice's Fashion",
			StoreDescription: "Trendy clothing and accessories",
		},
	}

	for _, user := range users {
		user.HashPassword("password123")
		if err := db.Create(user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
		} else {
			log.Printf("Created user: %s", user.Email)
		}
	}
}

func createSampleProducts(db *gorm.DB) {
	// Get sellers
	var sellers []models.User
	db.Where("role = ?", models.RoleSeller).Find(&sellers)
	
	if len(sellers) == 0 {
		log.Println("No sellers found, skipping product creation")
		return
	}

	products := []*models.Product{
		{
			Name:         "Wireless Bluetooth Headphones",
			Description:  "High-quality wireless headphones with noise cancellation and long battery life.",
			SKU:          "WHD-001",
			Price:        99.99,
			StockQuantity: 50,
			Category:     models.CategoryElectronics,
			Status:       models.ProductStatusActive,
			Visible:      true,
			SellerID:     sellers[0].ID,
		},
		{
			Name:         "Smartphone Case",
			Description:  "Durable protective case for smartphones with shock absorption.",
			SKU:          "SPC-001",
			Price:        19.99,
			StockQuantity: 100,
			Category:     models.CategoryElectronics,
			Status:       models.ProductStatusActive,
			Visible:      true,
			SellerID:     sellers[0].ID,
		},
		{
			Name:         "Cotton T-Shirt",
			Description:  "Comfortable 100% cotton t-shirt available in multiple colors.",
			SKU:          "CTS-001",
			Price:        24.99,
			StockQuantity: 200,
			Category:     models.CategoryClothing,
			Status:       models.ProductStatusActive,
			Visible:      true,
			SellerID:     sellers[1].ID,
		},
		{
			Name:         "Denim Jeans",
			Description:  "Classic blue denim jeans with comfortable fit.",
			SKU:          "DJ-001",
			Price:        59.99,
			StockQuantity: 75,
			Category:     models.CategoryClothing,
			Status:       models.ProductStatusActive,
			Visible:      true,
			SellerID:     sellers[1].ID,
		},
	}

	for _, product := range products {
		product.GenerateSlug()
		if err := db.Create(product).Error; err != nil {
			log.Printf("Failed to create product %s: %v", product.Name, err)
		} else {
			log.Printf("Created product: %s", product.Name)
		}
	}
}

func createSampleOrders(db *gorm.DB) {
	// Get customers and products
	var customers []models.User
	var products []models.Product
	
	db.Where("role = ?", models.RoleCustomer).Find(&customers)
	db.Where("status = ?", models.ProductStatusActive).Find(&products)
	
	if len(customers) == 0 || len(products) == 0 {
		log.Println("No customers or products found, skipping order creation")
		return
	}

	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 5; i++ {
		customer := customers[rand.Intn(len(customers))]
		product := products[rand.Intn(len(products))]
		
		order := &models.Order{
			CustomerID:         customer.ID,
			Status:            models.OrderStatusDelivered,
			PaymentStatus:     models.PaymentStatusPaid,
			PaymentMethod:     models.PaymentMethodCard,
			SubtotalAmount:    product.Price,
			TotalAmount:       product.Price + 5.99, // Add shipping
			ShippingAmount:    5.99,
			ShippingFirstName: customer.FirstName,
			ShippingLastName:  customer.LastName,
			ShippingEmail:     customer.Email,
			ShippingStreet:    "123 Main St",
			ShippingCity:      "Anytown",
			ShippingState:     "CA",
			ShippingCountry:   "US",
			ShippingPostalCode: "12345",
		}
		
		order.GenerateOrderNumber()
		
		if err := db.Create(order).Error; err != nil {
			log.Printf("Failed to create order: %v", err)
			continue
		}
		
		// Create order item
		orderItem := &models.OrderItem{
			OrderID:     order.ID,
			ProductID:   product.ID,
			Quantity:    1,
			UnitPrice:   product.Price,
			TotalPrice:  product.Price,
			ProductName: product.Name,
			ProductSKU:  product.SKU,
		}
		
		if err := db.Create(orderItem).Error; err != nil {
			log.Printf("Failed to create order item: %v", err)
		} else {
			log.Printf("Created order: %s", order.OrderNumber)
		}
	}
}

func createSampleReviews(db *gorm.DB) {
	// Get orders with delivered status
	var orders []models.Order
	db.Preload("Customer").Preload("OrderItems.Product").Where("status = ?", models.OrderStatusDelivered).Find(&orders)
	
	if len(orders) == 0 {
		log.Println("No delivered orders found, skipping review creation")
		return
	}

	reviews := []struct {
		Rating  int
		Title   string
		Comment string
	}{
		{5, "Excellent product!", "I'm very satisfied with this purchase. Great quality and fast shipping."},
		{4, "Good value for money", "Nice product, works as expected. Would recommend to others."},
		{5, "Love it!", "Amazing quality and exactly what I was looking for. Will buy again."},
		{3, "Decent product", "It's okay, does the job but nothing special."},
		{4, "Happy with purchase", "Good product, arrived on time and in perfect condition."},
	}

	rand.Seed(time.Now().UnixNano())
	
	for _, order := range orders {
		for _, item := range order.OrderItems {
			if rand.Float32() < 0.7 { // 70% chance of review
				reviewData := reviews[rand.Intn(len(reviews))]
				
				review := &models.Review{
					ProductID:  item.ProductID,
					UserID:     order.CustomerID,
					OrderID:    &order.ID,
					Rating:     reviewData.Rating,
					Title:      reviewData.Title,
					Comment:    reviewData.Comment,
					IsVerified: true,
					IsApproved: true,
				}
				
				if err := db.Create(review).Error; err != nil {
					log.Printf("Failed to create review: %v", err)
				} else {
					log.Printf("Created review for product ID %d", item.ProductID)
				}
			}
		}
	}
}
