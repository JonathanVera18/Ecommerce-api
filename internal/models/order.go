package models

import (
	"fmt"
	"time"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod represents payment methods
type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodPaypal PaymentMethod = "paypal"
	PaymentMethodBank   PaymentMethod = "bank_transfer"
	PaymentMethodCash   PaymentMethod = "cash_on_delivery"
)

// Order represents an order in the system
type Order struct {
	BaseModel
	OrderNumber string        `json:"order_number" gorm:"type:varchar(50);unique;not null"`
	CustomerID  uint          `json:"customer_id" gorm:"not null"`
	Customer    User          `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	
	// Order details
	Status        OrderStatus   `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	TotalAmount   float64       `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	SubtotalAmount float64      `json:"subtotal_amount" gorm:"type:decimal(10,2);not null"`
	TaxAmount     float64       `json:"tax_amount" gorm:"type:decimal(10,2);default:0"`
	ShippingAmount float64      `json:"shipping_amount" gorm:"type:decimal(10,2);default:0"`
	DiscountAmount float64      `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	
	// Payment information
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentMethod PaymentMethod `json:"payment_method" gorm:"type:varchar(20)"`
	PaymentID     *string       `json:"payment_id,omitempty" gorm:"type:varchar(255)"` // External payment ID
	PaidAt        *time.Time    `json:"paid_at,omitempty"`
	
	// Shipping information
	ShippingFirstName string  `json:"shipping_first_name" gorm:"type:varchar(100);not null"`
	ShippingLastName  string  `json:"shipping_last_name" gorm:"type:varchar(100);not null"`
	ShippingEmail     string  `json:"shipping_email" gorm:"type:varchar(255)"`
	ShippingPhone     *string `json:"shipping_phone,omitempty" gorm:"type:varchar(20)"`
	ShippingStreet    string  `json:"shipping_street" gorm:"type:varchar(255);not null"`
	ShippingCity      string  `json:"shipping_city" gorm:"type:varchar(100);not null"`
	ShippingState     string  `json:"shipping_state" gorm:"type:varchar(100);not null"`
	ShippingCountry   string  `json:"shipping_country" gorm:"type:varchar(100);not null"`
	ShippingPostalCode string `json:"shipping_postal_code" gorm:"type:varchar(20);not null"`
	
	// Billing information (optional, can be same as shipping)
	BillingFirstName   *string `json:"billing_first_name,omitempty" gorm:"type:varchar(100)"`
	BillingLastName    *string `json:"billing_last_name,omitempty" gorm:"type:varchar(100)"`
	BillingEmail       *string `json:"billing_email,omitempty" gorm:"type:varchar(255)"`
	BillingPhone       *string `json:"billing_phone,omitempty" gorm:"type:varchar(20)"`
	BillingStreet      *string `json:"billing_street,omitempty" gorm:"type:varchar(255)"`
	BillingCity        *string `json:"billing_city,omitempty" gorm:"type:varchar(100)"`
	BillingState       *string `json:"billing_state,omitempty" gorm:"type:varchar(100)"`
	BillingCountry     *string `json:"billing_country,omitempty" gorm:"type:varchar(100)"`
	BillingPostalCode  *string `json:"billing_postal_code,omitempty" gorm:"type:varchar(20)"`
	
	// Tracking information
	TrackingNumber *string    `json:"tracking_number,omitempty" gorm:"type:varchar(100)"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	
	// Additional information
	Notes        *string `json:"notes,omitempty" gorm:"type:text"`
	InternalNotes *string `json:"internal_notes,omitempty" gorm:"type:text"` // Admin/staff notes
	
	// Relationships
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	
	// Computed fields
	ItemCount int `json:"item_count" gorm:"-"`
}

// OrderItem represents items in an order
type OrderItem struct {
	BaseModel
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	
	Quantity  int     `json:"quantity" gorm:"not null" validate:"min=1"`
	UnitPrice float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	TotalPrice float64 `json:"total_price" gorm:"type:decimal(10,2);not null"`
	
	// Product snapshot (to preserve product details at time of order)
	ProductName        string  `json:"product_name" gorm:"type:varchar(255);not null"`
	ProductSKU         string  `json:"product_sku" gorm:"type:varchar(100);not null"`
	ProductDescription *string `json:"product_description,omitempty" gorm:"type:text"`
	ProductImage       *string `json:"product_image,omitempty" gorm:"type:varchar(500)"`
}

// Cart represents a shopping cart (temporary before order)
type Cart struct {
	BaseModel
	CustomerID uint       `json:"customer_id" gorm:"not null;unique"`
	Customer   User       `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	CartItems  []CartItem `json:"cart_items,omitempty" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	
	// Computed fields
	TotalAmount float64 `json:"total_amount" gorm:"-"`
	ItemCount   int     `json:"item_count" gorm:"-"`
}

// CartItem represents items in a cart
type CartItem struct {
	BaseModel
	CartID    uint    `json:"cart_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null" validate:"min=1"`
}

// CartAddRequest represents the request to add item to cart
type CartAddRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

// CartUpdateRequest represents the request to update cart item
type CartUpdateRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// CartResponse represents the cart response
type CartResponse struct {
	ID          uint               `json:"id"`
	CustomerID  uint               `json:"customer_id"`
	Items       []CartItemResponse `json:"items"`
	TotalAmount float64            `json:"total_amount"`
	ItemCount   int                `json:"item_count"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// CartItemResponse represents the cart item response
type CartItemResponse struct {
	ID        uint            `json:"id"`
	CartID    uint            `json:"cart_id"`
	ProductID uint            `json:"product_id"`
	Product   ProductResponse `json:"product"`
	Quantity  int             `json:"quantity"`
	Subtotal  float64         `json:"subtotal"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// OrderCreateRequest represents the request to create an order
type OrderCreateRequest struct {
	PaymentMethod PaymentMethod `json:"payment_method" validate:"required"`
	
	// Shipping information
	ShippingFirstName  string  `json:"shipping_first_name" validate:"required,min=2,max=100"`
	ShippingLastName   string  `json:"shipping_last_name" validate:"required,min=2,max=100"`
	ShippingEmail      string  `json:"shipping_email" validate:"required,email"`
	ShippingPhone      *string `json:"shipping_phone,omitempty" validate:"omitempty,e164"`
	ShippingStreet     string  `json:"shipping_street" validate:"required,min=5,max=255"`
	ShippingCity       string  `json:"shipping_city" validate:"required,min=2,max=100"`
	ShippingState      string  `json:"shipping_state" validate:"required,min=2,max=100"`
	ShippingCountry    string  `json:"shipping_country" validate:"required,min=2,max=100"`
	ShippingPostalCode string  `json:"shipping_postal_code" validate:"required,min=3,max=20"`
	
	// Billing information (optional)
	BillingFirstName   *string `json:"billing_first_name,omitempty" validate:"omitempty,min=2,max=100"`
	BillingLastName    *string `json:"billing_last_name,omitempty" validate:"omitempty,min=2,max=100"`
	BillingEmail       *string `json:"billing_email,omitempty" validate:"omitempty,email"`
	BillingPhone       *string `json:"billing_phone,omitempty" validate:"omitempty,e164"`
	BillingStreet      *string `json:"billing_street,omitempty" validate:"omitempty,min=5,max=255"`
	BillingCity        *string `json:"billing_city,omitempty" validate:"omitempty,min=2,max=100"`
	BillingState       *string `json:"billing_state,omitempty" validate:"omitempty,min=2,max=100"`
	BillingCountry     *string `json:"billing_country,omitempty" validate:"omitempty,min=2,max=100"`
	BillingPostalCode  *string `json:"billing_postal_code,omitempty" validate:"omitempty,min=3,max=20"`
	
	Notes *string `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// OrderUpdateRequest represents the request to update an order
type OrderUpdateRequest struct {
	Status         *OrderStatus   `json:"status,omitempty"`
	TrackingNumber *string        `json:"tracking_number,omitempty" validate:"omitempty,max=100"`
	InternalNotes  *string        `json:"internal_notes,omitempty" validate:"omitempty,max=1000"`
}

// OrderListRequest represents the request to list orders with filters
type OrderListRequest struct {
	Page          int            `query:"page" validate:"min=1"`
	Limit         int            `query:"limit" validate:"min=1,max=100"`
	Status        *OrderStatus   `query:"status"`
	PaymentStatus *PaymentStatus `query:"payment_status"`
	CustomerID    *uint          `query:"customer_id"`
	DateFrom      *time.Time     `query:"date_from"`
	DateTo        *time.Time     `query:"date_to"`
	MinAmount     *float64       `query:"min_amount" validate:"omitempty,min=0"`
	MaxAmount     *float64       `query:"max_amount" validate:"omitempty,min=0"`
	Search        string         `query:"search"`
	SortBy        string         `query:"sort_by" validate:"omitempty,oneof=created_at updated_at total_amount order_number"`
	SortOrder     string         `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// CartItemRequest represents the request to add/update cart items
type CartItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

// OrderStatsResponse represents order statistics
type OrderStatsResponse struct {
	TotalOrders     int64   `json:"total_orders"`
	PendingOrders   int64   `json:"pending_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID       uint          `json:"order_id" validate:"required"`
	PaymentMethod PaymentMethod `json:"payment_method" validate:"required"`
	Amount        float64       `json:"amount" validate:"required,min=0.01"`
	Currency      string        `json:"currency" validate:"required,len=3"`
	
	// Stripe specific fields
	PaymentMethodID *string `json:"payment_method_id,omitempty"` // Stripe payment method ID
	
	// Return URLs
	SuccessURL string `json:"success_url" validate:"required,url"`
	CancelURL  string `json:"cancel_url" validate:"required,url"`
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	Items           []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
	ShippingAddress string             `json:"shipping_address" validate:"required"`
	PaymentMethod   PaymentMethod      `json:"payment_method" validate:"required"`
}

// OrderItemRequest represents an order item in a request
type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

// UpdateOrderStatusRequest represents the request to update order status
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
}

// PaymentProcessRequest represents a payment processing request
type PaymentProcessRequest struct {
	Token string `json:"token" validate:"required"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Amount        float64 `json:"amount"`
}

// OrderAnalytics represents order analytics data
type OrderAnalytics struct {
	TotalRevenue     float64 `json:"total_revenue"`
	TotalOrders      int64   `json:"total_orders"`
	PendingOrders    int64   `json:"pending_orders"`
	ConfirmedOrders  int64   `json:"confirmed_orders"`
	ShippedOrders    int64   `json:"shipped_orders"`
	DeliveredOrders  int64   `json:"delivered_orders"`
	CancelledOrders  int64   `json:"cancelled_orders"`
}

// GenerateOrderNumber generates a unique order number
func (o *Order) GenerateOrderNumber() {
	// Format: ORD-YYYYMMDD-HHMMSS-XXX (XXX is random)
	now := time.Now()
	o.OrderNumber = fmt.Sprintf("ORD-%s-%03d", 
		now.Format("20060102-150405"), 
		o.ID%1000)
}

// CalculateTotals calculates order totals based on order items
func (o *Order) CalculateTotals() {
	o.SubtotalAmount = 0
	o.ItemCount = 0
	
	for _, item := range o.OrderItems {
		o.SubtotalAmount += item.TotalPrice
		o.ItemCount += item.Quantity
	}
	
	// Calculate total (subtotal + tax + shipping - discount)
	o.TotalAmount = o.SubtotalAmount + o.TaxAmount + o.ShippingAmount - o.DiscountAmount
}

// CanCancel checks if the order can be cancelled
func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed
}

// CanRefund checks if the order can be refunded
func (o *Order) CanRefund() bool {
	return o.PaymentStatus == PaymentStatusPaid && 
		(o.Status == OrderStatusDelivered || o.Status == OrderStatusShipped)
}

// CanShip checks if the order can be shipped
func (o *Order) CanShip() bool {
	return o.Status == OrderStatusConfirmed || o.Status == OrderStatusProcessing
}

// IsCompleted checks if the order is completed
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusDelivered
}

// IsCancelled checks if the order is cancelled
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// GetShippingAddress returns formatted shipping address
func (o *Order) GetShippingAddress() string {
	return fmt.Sprintf("%s %s\n%s\n%s, %s %s\n%s",
		o.ShippingFirstName, o.ShippingLastName,
		o.ShippingStreet,
		o.ShippingCity, o.ShippingState, o.ShippingPostalCode,
		o.ShippingCountry)
}

// CalculateTotal calculates total price for order item
func (oi *OrderItem) CalculateTotal() {
	oi.TotalPrice = oi.UnitPrice * float64(oi.Quantity)
}

// UpdateFromProduct updates order item fields from product
func (oi *OrderItem) UpdateFromProduct(product *Product) {
	oi.ProductName = product.Name
	oi.ProductSKU = product.SKU
	oi.ProductDescription = &product.Description
	// Get the primary image URL directly as a string
	primaryImage := product.GetPrimaryImage()
	oi.ProductImage = &primaryImage
	oi.UnitPrice = product.Price
	oi.CalculateTotal()
}

// CalculateTotals calculates cart totals
func (c *Cart) CalculateTotals() {
	c.TotalAmount = 0
	c.ItemCount = 0
	
	for _, item := range c.CartItems {
		if item.Product.ID != 0 {
			c.TotalAmount += item.Product.Price * float64(item.Quantity)
		}
		c.ItemCount += item.Quantity
	}
}

// ToResponse converts Cart to CartResponse
func (c *Cart) ToResponse() CartResponse {
	resp := CartResponse{
		ID:         c.ID,
		CustomerID: c.CustomerID,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
	
	totalAmount := 0.0
	itemCount := 0
	
	for _, item := range c.CartItems {
		itemResp := CartItemResponse{
			ID:        item.ID,
			CartID:    item.CartID,
			ProductID: item.ProductID,
			Product:   item.Product.ToResponse(),
			Quantity:  item.Quantity,
			Subtotal:  item.Product.Price * float64(item.Quantity),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
		resp.Items = append(resp.Items, itemResp)
		totalAmount += itemResp.Subtotal
		itemCount += item.Quantity
	}
	
	resp.TotalAmount = totalAmount
	resp.ItemCount = itemCount
	
	return resp
}