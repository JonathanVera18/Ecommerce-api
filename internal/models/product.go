package models

import (
	"fmt"
	"strings"
	"time"
)

// ProductStatus represents product status
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusInactive  ProductStatus = "inactive"
	ProductStatusDeleted   ProductStatus = "deleted"
)

// ProductCategory represents product categories
type ProductCategory string

const (
	CategoryElectronics ProductCategory = "electronics"
	CategoryClothing    ProductCategory = "clothing"
	CategoryBooks       ProductCategory = "books"
	CategoryHome        ProductCategory = "home"
	CategorySports      ProductCategory = "sports"
	CategoryToys        ProductCategory = "toys"
	CategoryBeauty      ProductCategory = "beauty"
	CategoryFood        ProductCategory = "food"
	CategoryOther       ProductCategory = "other"
)

// Product represents a product in the system
type Product struct {
	BaseModel
	Name         string          `json:"name" gorm:"type:varchar(255);not null" validate:"required,min=3,max=255"`
	Description  string          `json:"description" gorm:"type:text" validate:"required,min=10"`
	ShortDesc    *string         `json:"short_description,omitempty" gorm:"type:varchar(500)"`
	SKU          string          `json:"sku" gorm:"type:varchar(100);unique;not null" validate:"required"`
	Price        float64         `json:"price" gorm:"type:decimal(10,2);not null" validate:"required,min=0"`
	ComparePrice *float64        `json:"compare_price,omitempty" gorm:"type:decimal(10,2)" validate:"omitempty,gtfield=Price"`
	CostPrice    *float64        `json:"cost_price,omitempty" gorm:"type:decimal(10,2)" validate:"omitempty,min=0"`
	
	// Inventory - simplified for compatibility
	Stock       int  `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	StockQuantity    int  `json:"stock_quantity" gorm:"not null;default:0" validate:"min=0"`
	LowStockLevel    int  `json:"low_stock_level" gorm:"default:10" validate:"min=0"`
	TrackInventory   bool `json:"track_inventory" gorm:"default:true"`
	AllowBackorders  bool `json:"allow_backorders" gorm:"default:false"`
	
	// Organization
	Category   string `json:"category" gorm:"type:varchar(50);not null" validate:"required"`
	CategoryID *uint  `json:"category_id,omitempty" gorm:"index"`
	Tags       string `json:"tags,omitempty" gorm:"type:varchar(1000)"` // Comma-separated tags
	Brand      *string `json:"brand,omitempty" gorm:"type:varchar(100)"`
	
	// Physical properties
	Weight     *float64 `json:"weight,omitempty" gorm:"type:decimal(8,3)" validate:"omitempty,min=0"`
	Length     *float64 `json:"length,omitempty" gorm:"type:decimal(8,2)" validate:"omitempty,min=0"`
	Width      *float64 `json:"width,omitempty" gorm:"type:decimal(8,2)" validate:"omitempty,min=0"`
	Height     *float64 `json:"height,omitempty" gorm:"type:decimal(8,2)" validate:"omitempty,min=0"`
	
	// SEO
	MetaTitle       *string `json:"meta_title,omitempty" gorm:"type:varchar(255)"`
	MetaDescription *string `json:"meta_description,omitempty" gorm:"type:varchar(500)"`
	Slug            string  `json:"slug" gorm:"type:varchar(255);unique;not null"`
	
	// Status and visibility - simplified for compatibility
	IsActive  bool          `json:"is_active" gorm:"default:true"`
	Status    ProductStatus `json:"status" gorm:"type:varchar(20);not null;default:'draft'" validate:"required"`
	Featured  bool          `json:"featured" gorm:"default:false"`
	Visible   bool          `json:"visible" gorm:"default:true"`
	
	// Images - simplified for compatibility
	Images []string `json:"images,omitempty" gorm:"-"`
	ProductImages []ProductImage `json:"product_images,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	
	// Seller
	SellerID uint `json:"seller_id" gorm:"not null"`
	Seller   User `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
	
	// Analytics
	ViewCount int `json:"view_count" gorm:"default:0"`
	
	// Relationships
	OrderItems []OrderItem `json:"-" gorm:"foreignKey:ProductID"`
	Reviews    []Review    `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
	
	// Computed fields (not stored in DB)
	AverageRating float64 `json:"average_rating" gorm:"column:average_rating;default:0"`
	ReviewCount   int     `json:"review_count" gorm:"column:review_count;default:0"`
	IsLowStock    bool    `json:"is_low_stock" gorm:"-"`
	IsInStock     bool    `json:"is_in_stock" gorm:"-"`
}

// ProductImage represents product images
type ProductImage struct {
	BaseModel
	ProductID uint   `json:"product_id" gorm:"not null"`
	URL       string `json:"url" gorm:"type:varchar(500);not null" validate:"required,url"`
	AltText   string `json:"alt_text" gorm:"type:varchar(255)" validate:"max=255"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
	IsPrimary bool   `json:"is_primary" gorm:"default:false"`
}

// ProductCreateRequest represents the request to create a product
type ProductCreateRequest struct {
	Name         string          `json:"name" validate:"required,min=3,max=255"`
	Description  string          `json:"description" validate:"required,min=10"`
	ShortDesc    *string         `json:"short_description,omitempty" validate:"omitempty,max=500"`
	SKU          string          `json:"sku" validate:"required"`
	Price        float64         `json:"price" validate:"required,min=0"`
	ComparePrice *float64        `json:"compare_price,omitempty" validate:"omitempty,gtfield=Price"`
	CostPrice    *float64        `json:"cost_price,omitempty" validate:"omitempty,min=0"`
	
	StockQuantity   int             `json:"stock_quantity" validate:"min=0"`
	LowStockLevel   int             `json:"low_stock_level" validate:"min=0"`
	TrackInventory  bool            `json:"track_inventory"`
	AllowBackorders bool            `json:"allow_backorders"`
	
	Category ProductCategory `json:"category" validate:"required"`
	Tags     []string        `json:"tags,omitempty"`
	Brand    *string         `json:"brand,omitempty" validate:"omitempty,max=100"`
	
	Weight *float64 `json:"weight,omitempty" validate:"omitempty,min=0"`
	Length *float64 `json:"length,omitempty" validate:"omitempty,min=0"`
	Width  *float64 `json:"width,omitempty" validate:"omitempty,min=0"`
	Height *float64 `json:"height,omitempty" validate:"omitempty,min=0"`
	
	MetaTitle       *string `json:"meta_title,omitempty" validate:"omitempty,max=255"`
	MetaDescription *string `json:"meta_description,omitempty" validate:"omitempty,max=500"`
	
	Status   ProductStatus `json:"status" validate:"required"`
	Featured bool          `json:"featured"`
	Visible  bool          `json:"visible"`
	
	Images []ProductImageRequest `json:"images,omitempty"`
}

// ProductUpdateRequest represents the request to update a product
type ProductUpdateRequest struct {
	Name         *string         `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description  *string         `json:"description,omitempty" validate:"omitempty,min=10"`
	ShortDesc    *string         `json:"short_description,omitempty" validate:"omitempty,max=500"`
	Price        *float64        `json:"price,omitempty" validate:"omitempty,min=0"`
	ComparePrice *float64        `json:"compare_price,omitempty"`
	CostPrice    *float64        `json:"cost_price,omitempty" validate:"omitempty,min=0"`
	
	StockQuantity   *int  `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	LowStockLevel   *int  `json:"low_stock_level,omitempty" validate:"omitempty,min=0"`
	TrackInventory  *bool `json:"track_inventory,omitempty"`
	AllowBackorders *bool `json:"allow_backorders,omitempty"`
	
	Category *ProductCategory `json:"category,omitempty"`
	Tags     []string         `json:"tags,omitempty"`
	Brand    *string          `json:"brand,omitempty" validate:"omitempty,max=100"`
	
	Weight *float64 `json:"weight,omitempty" validate:"omitempty,min=0"`
	Length *float64 `json:"length,omitempty" validate:"omitempty,min=0"`
	Width  *float64 `json:"width,omitempty" validate:"omitempty,min=0"`
	Height *float64 `json:"height,omitempty" validate:"omitempty,min=0"`
	
	MetaTitle       *string `json:"meta_title,omitempty" validate:"omitempty,max=255"`
	MetaDescription *string `json:"meta_description,omitempty" validate:"omitempty,max=500"`
	
	Status   *ProductStatus `json:"status,omitempty"`
	Featured *bool          `json:"featured,omitempty"`
	Visible  *bool          `json:"visible,omitempty"`
}

// ProductImageRequest represents the request to add/update product images
type ProductImageRequest struct {
	URL       string `json:"url" validate:"required,url"`
	AltText   string `json:"alt_text" validate:"max=255"`
	SortOrder int    `json:"sort_order"`
	IsPrimary bool   `json:"is_primary"`
}

// ProductListRequest represents the request to list products with filters
type ProductListRequest struct {
	Page         int               `query:"page" validate:"min=1"`
	Limit        int               `query:"limit" validate:"min=1,max=100"`
	Category     *ProductCategory  `query:"category"`
	Status       *ProductStatus    `query:"status"`
	SellerID     *uint             `query:"seller_id"`
	MinPrice     *float64          `query:"min_price" validate:"omitempty,min=0"`
	MaxPrice     *float64          `query:"max_price" validate:"omitempty,min=0"`
	InStock      *bool             `query:"in_stock"`
	Featured     *bool             `query:"featured"`
	Search       string            `query:"search"`
	SortBy       string            `query:"sort_by" validate:"omitempty,oneof=name price created_at updated_at view_count rating"`
	SortOrder    string            `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// ProductStatsResponse represents product statistics
type ProductStatsResponse struct {
	TotalProducts    int64   `json:"total_products"`
	ActiveProducts   int64   `json:"active_products"`
	LowStockProducts int64   `json:"low_stock_products"`
	OutOfStock       int64   `json:"out_of_stock"`
	TotalValue       float64 `json:"total_value"`
	AveragePrice     float64 `json:"average_price"`
}

// Request models
type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required,min=3,max=255"`
	Description string   `json:"description" validate:"required,min=10"`
	Price       float64  `json:"price" validate:"required,min=0"`
	Stock       int      `json:"stock" validate:"min=0"`
	Category    string   `json:"category" validate:"required"`
	Images      []string `json:"images,omitempty"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=10"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	Category    *string  `json:"category,omitempty"`
	Images      []string `json:"images,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

type GetProductsRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Category string `json:"category,omitempty"`
	Search   string `json:"search,omitempty"`
	SellerID *uint  `json:"seller_id,omitempty"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"min=0"`
}

// Response models
type ProductListResponse struct {
	Products []*Product `json:"products"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	Limit    int        `json:"limit"`
}

// ProductResponse represents the product response
type ProductResponse struct {
	ID              uint                    `json:"id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	ShortDesc       *string                 `json:"short_description,omitempty"`
	SKU             string                  `json:"sku"`
	Price           float64                 `json:"price"`
	ComparePrice    *float64                `json:"compare_price,omitempty"`
	CostPrice       *float64                `json:"cost_price,omitempty"`
	Stock           int                     `json:"stock"`
	StockQuantity   int                     `json:"stock_quantity"`
	LowStockLevel   int                     `json:"low_stock_level"`
	TrackInventory  bool                    `json:"track_inventory"`
	AllowBackorders bool                    `json:"allow_backorders"`
	Category        string                  `json:"category"`
	CategoryID      *uint                   `json:"category_id,omitempty"`
	Tags            []string                `json:"tags,omitempty"`
	Brand           *string                 `json:"brand,omitempty"`
	Weight          *float64                `json:"weight,omitempty"`
	Length          *float64                `json:"length,omitempty"`
	Width           *float64                `json:"width,omitempty"`
	Height          *float64                `json:"height,omitempty"`
	MetaTitle       *string                 `json:"meta_title,omitempty"`
	MetaDescription *string                 `json:"meta_description,omitempty"`
	Slug            string                  `json:"slug"`
	IsActive        bool                    `json:"is_active"`
	Status          ProductStatus           `json:"status"`
	Featured        bool                    `json:"featured"`
	Visible         bool                    `json:"visible"`
	Images          []string                `json:"images,omitempty"`
	ProductImages   []ProductImageResponse  `json:"product_images,omitempty"`
	SellerID        uint                    `json:"seller_id"`
	ViewCount       int                     `json:"view_count"`
	AverageRating   float64                 `json:"average_rating"`
	ReviewCount     int                     `json:"review_count"`
	IsLowStock      bool                    `json:"is_low_stock"`
	IsInStock       bool                    `json:"is_in_stock"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	
	// Computed fields
	DiscountPercent float64 `json:"discount_percent"`
	PrimaryImage    string  `json:"primary_image"`
}

// ProductImageResponse represents product image response
type ProductImageResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	p.UpdateComputedFields()
	
	resp := ProductResponse{
		ID:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		ShortDesc:       p.ShortDesc,
		SKU:             p.SKU,
		Price:           p.Price,
		ComparePrice:    p.ComparePrice,
		CostPrice:       p.CostPrice,
		Stock:           p.Stock,
		StockQuantity:   p.StockQuantity,
		LowStockLevel:   p.LowStockLevel,
		TrackInventory:  p.TrackInventory,
		AllowBackorders: p.AllowBackorders,
		Category:        p.Category,
		CategoryID:      p.CategoryID,
		Tags:            p.GetTagsList(),
		Brand:           p.Brand,
		Weight:          p.Weight,
		Length:          p.Length,
		Width:           p.Width,
		Height:          p.Height,
		MetaTitle:       p.MetaTitle,
		MetaDescription: p.MetaDescription,
		Slug:            p.Slug,
		IsActive:        p.IsActive,
		Status:          p.Status,
		Featured:        p.Featured,
		Visible:         p.Visible,
		Images:          p.Images,
		SellerID:        p.SellerID,
		ViewCount:       p.ViewCount,
		AverageRating:   p.AverageRating,
		ReviewCount:     p.ReviewCount,
		IsLowStock:      p.IsLowStock,
		IsInStock:       p.IsInStock,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		DiscountPercent: p.CalculateDiscount(),
		PrimaryImage:    p.GetPrimaryImage(),
	}
	
	// Convert product images
	if p.ProductImages != nil {
		for _, img := range p.ProductImages {
			resp.ProductImages = append(resp.ProductImages, ProductImageResponse{
				ID:        img.ID,
				ProductID: img.ProductID,
				URL:       img.URL,
				AltText:   img.AltText,
				SortOrder: img.SortOrder,
				IsPrimary: img.IsPrimary,
				CreatedAt: img.CreatedAt,
				UpdatedAt: img.UpdatedAt,
			})
		}
	}
	
	return resp
}

// GetTagsList returns tags as a slice
func (p *Product) GetTagsList() []string {
	if p.Tags == "" {
		return []string{}
	}
	tags := strings.Split(p.Tags, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}

// SetTagsList sets tags from a slice
func (p *Product) SetTagsList(tags []string) {
	if len(tags) == 0 {
		p.Tags = ""
		return
	}
	p.Tags = strings.Join(tags, ",")
}

// UpdateComputedFields updates computed fields
func (p *Product) UpdateComputedFields() {
	p.IsLowStock = p.TrackInventory && p.StockQuantity <= p.LowStockLevel
	p.IsInStock = !p.TrackInventory || p.StockQuantity > 0
}

// GenerateSlug generates a URL-friendly slug from the product name
func (p *Product) GenerateSlug() {
	slug := strings.ToLower(p.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	p.Slug = result.String()
}

// CanOrder checks if the product can be ordered
func (p *Product) CanOrder(quantity int) bool {
	if p.Status != ProductStatusActive || !p.Visible {
		return false
	}
	
	if !p.TrackInventory {
		return true
	}
	
	if p.StockQuantity >= quantity {
		return true
	}
	
	return p.AllowBackorders
}

// CalculateDiscount calculates discount percentage if compare price is set
func (p *Product) CalculateDiscount() float64 {
	if p.ComparePrice == nil || *p.ComparePrice <= p.Price {
		return 0
	}
	return ((*p.ComparePrice - p.Price) / *p.ComparePrice) * 100
}

// GetPrimaryImage returns the primary image URL
func (p *Product) GetPrimaryImage() string {
	// First check ProductImages (the structured image objects)
	if len(p.ProductImages) > 0 {
		for _, img := range p.ProductImages {
			if img.IsPrimary {
				return img.URL
			}
		}
		// Return first image if no primary found
		return p.ProductImages[0].URL
	}
	
	// Fallback to simple Images slice
	if len(p.Images) > 0 {
		return p.Images[0]
	}
	
	return ""
}

// HasSufficientStock checks if there's sufficient stock for the given quantity
func (p *Product) HasSufficientStock(quantity int) bool {
	if !p.TrackInventory {
		return true
	}
	return p.StockQuantity >= quantity
}

// ReserveStock reduces stock quantity (for order processing)
func (p *Product) ReserveStock(quantity int) error {
	if !p.TrackInventory {
		return nil
	}
	
	if p.StockQuantity < quantity && !p.AllowBackorders {
		return fmt.Errorf("insufficient stock: available %d, requested %d", p.StockQuantity, quantity)
	}
	
	p.StockQuantity -= quantity
	return nil
}

// RestoreStock increases stock quantity (for order cancellation)
func (p *Product) RestoreStock(quantity int) {
	if !p.TrackInventory {
		return
	}
	p.StockQuantity += quantity
}