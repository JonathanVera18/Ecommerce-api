package models

import (
	"time"
)

// Category represents a product category
type Category struct {
	BaseModel
	Name        string  `json:"name" gorm:"type:varchar(100);not null;unique" validate:"required,min=2,max=100"`
	Slug        string  `json:"slug" gorm:"type:varchar(100);not null;unique" validate:"required"`
	Description *string `json:"description,omitempty" gorm:"type:text"`
	ImageURL    *string `json:"image_url,omitempty" gorm:"type:varchar(500)" validate:"omitempty,url"`
	ParentID    *uint   `json:"parent_id,omitempty" gorm:"index"`
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	SortOrder   int     `json:"sort_order" gorm:"default:0"`
	
	// Relationships
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []Product  `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	
	// Computed fields
	ProductCount int `json:"product_count" gorm:"-"`
}

// CategoryCreateRequest represents the request to create a category
type CategoryCreateRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty" validate:"omitempty,url"`
	ParentID    *uint   `json:"parent_id,omitempty"`
	IsActive    bool    `json:"is_active"`
	SortOrder   int     `json:"sort_order"`
}

// CategoryUpdateRequest represents the request to update a category
type CategoryUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty" validate:"omitempty,url"`
	ParentID    *uint   `json:"parent_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

// CategoryResponse represents the category response
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	ParentID    *uint     `json:"parent_id,omitempty"`
	IsActive    bool      `json:"is_active"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Parent       *CategoryResponse   `json:"parent,omitempty"`
	Children     []CategoryResponse  `json:"children,omitempty"`
	ProductCount int                 `json:"product_count"`
}

// ToResponse converts Category to CategoryResponse
func (c *Category) ToResponse() CategoryResponse {
	resp := CategoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		ImageURL:     c.ImageURL,
		ParentID:     c.ParentID,
		IsActive:     c.IsActive,
		SortOrder:    c.SortOrder,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		ProductCount: c.ProductCount,
	}
	
	if c.Parent != nil {
		parentResp := c.Parent.ToResponse()
		resp.Parent = &parentResp
	}
	
	if c.Children != nil {
		for _, child := range c.Children {
			resp.Children = append(resp.Children, child.ToResponse())
		}
	}
	
	return resp
}