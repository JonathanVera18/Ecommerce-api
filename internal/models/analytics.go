package models

import (
	"time"
)

// Dashboard statistics
type DashboardStats struct {
	TotalRevenue     float64 `json:"total_revenue"`
	TotalOrders      int64   `json:"total_orders"`
	PendingOrders    int64   `json:"pending_orders"`
	ConfirmedOrders  int64   `json:"confirmed_orders"`
	ShippedOrders    int64   `json:"shipped_orders"`
	DeliveredOrders  int64   `json:"delivered_orders"`
	CancelledOrders  int64   `json:"cancelled_orders"`
	TotalUsers       int64   `json:"total_users"`
	ActiveUsers      int64   `json:"active_users"`
	TotalProducts    int64   `json:"total_products"`
	LowStockProducts int     `json:"low_stock_products"`
}

// Sales analytics
type SalesAnalytics struct {
	Period           string    `json:"period"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	TotalRevenue     float64   `json:"total_revenue"`
	TotalOrders      int64     `json:"total_orders"`
	AverageOrderValue float64  `json:"average_order_value"`
	// You can add daily/weekly/monthly breakdowns here
	DailyBreakdown   []DailySales   `json:"daily_breakdown,omitempty"`
	WeeklyBreakdown  []WeeklySales  `json:"weekly_breakdown,omitempty"`
	MonthlyBreakdown []MonthlySales `json:"monthly_breakdown,omitempty"`
}

type DailySales struct {
	Date    time.Time `json:"date"`
	Revenue float64   `json:"revenue"`
	Orders  int64     `json:"orders"`
}

type WeeklySales struct {
	Week    int     `json:"week"`
	Year    int     `json:"year"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

type MonthlySales struct {
	Month   int     `json:"month"`
	Year    int     `json:"year"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

// User analytics
type UserAnalytics struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	NewUsers      int64 `json:"new_users"`
	CustomerCount int64 `json:"customer_count"`
	SellerCount   int64 `json:"seller_count"`
	AdminCount    int64 `json:"admin_count"`
}

// Product analytics
type ProductAnalytics struct {
	TotalProducts      int64       `json:"total_products"`
	LowStockProducts   int         `json:"low_stock_products"`
	OutOfStockProducts int         `json:"out_of_stock_products"`
	TopRatedProducts   []*Product  `json:"top_rated_products"`
}

// Review analytics
type ReviewAnalytics struct {
	TotalReviews   int64     `json:"total_reviews"`
	AverageRating  float64   `json:"average_rating"`
	RecentReviews  []*Review `json:"recent_reviews"`
	TopReviews     []*Review `json:"top_reviews"`
}

// System health
type SystemHealth struct {
	Status         string        `json:"status"`
	DatabaseStatus string        `json:"database_status"`
	RedisStatus    string        `json:"redis_status"`
	LastChecked    time.Time     `json:"last_checked"`
	Uptime         time.Duration `json:"uptime"`
}

// Admin user management request
type AdminUserUpdateRequest struct {
	Role      *UserRole `json:"role,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
	IsVerified *bool    `json:"is_verified,omitempty"`
}
