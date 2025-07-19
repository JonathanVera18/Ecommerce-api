package models

import (
	"time"
)

// Wishlist represents a user's wishlist
type Wishlist struct {
	BaseModel
	UserID    uint `json:"user_id" gorm:"not null;index"`
	ProductID uint `json:"product_id" gorm:"not null;index"`
	
	// Relationships
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// WishlistAddRequest represents the request to add item to wishlist
type WishlistAddRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
}

// WishlistResponse represents the wishlist response
type WishlistResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// Product information
	Product *ProductResponse `json:"product,omitempty"`
}

// ToResponse converts Wishlist to WishlistResponse
func (w *Wishlist) ToResponse() WishlistResponse {
	resp := WishlistResponse{
		ID:        w.ID,
		UserID:    w.UserID,
		ProductID: w.ProductID,
		CreatedAt: w.CreatedAt,
	}
	
	if w.Product.ID != 0 {
		productResp := w.Product.ToResponse()
		resp.Product = &productResp
	}
	
	return resp
}

// WishlistItemsResponse represents the wishlist items response
type WishlistItemsResponse struct {
	Items []WishlistResponse `json:"items"`
	Total int                `json:"total"`
}
