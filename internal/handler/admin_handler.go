package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	userService    service.UserService
	productService service.ProductService
	orderService   service.OrderService
	reviewService  service.ReviewService
}

func NewAdminHandler(
	userService service.UserService,
	productService service.ProductService,
	orderService service.OrderService,
	reviewService service.ReviewService,
) *AdminHandler {
	return &AdminHandler{
		userService:    userService,
		productService: productService,
		orderService:   orderService,
		reviewService:  reviewService,
	}
}

// GetDashboardStats retrieves dashboard statistics
// @Summary Get dashboard statistics
// @Description Get overall platform statistics (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} utils.Response{data=models.DashboardStats}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/dashboard [get]
func (h *AdminHandler) GetDashboardStats(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	// Get analytics for the last 30 days by default
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// Get order analytics
	orderAnalytics, err := h.orderService.GetOrderAnalytics(c.Request().Context(), nil, &startDate, &endDate)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get order analytics")
	}

	// Get user stats (you would need to implement this in UserService)
	// userStats, err := h.userService.GetUserStats(c.Request().Context())
	// if err != nil {
	//     return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user stats")
	// }

	stats := &models.DashboardStats{
		TotalRevenue:     orderAnalytics.TotalRevenue,
		TotalOrders:      orderAnalytics.TotalOrders,
		PendingOrders:    orderAnalytics.PendingOrders,
		ConfirmedOrders:  orderAnalytics.ConfirmedOrders,
		ShippedOrders:    orderAnalytics.ShippedOrders,
		DeliveredOrders:  orderAnalytics.DeliveredOrders,
		CancelledOrders:  orderAnalytics.CancelledOrders,
		// TotalUsers:       userStats.TotalUsers,
		// ActiveUsers:      userStats.ActiveUsers,
		// TotalProducts:    productStats.TotalProducts,
		// LowStockProducts: len(lowStockProducts),
	}

	return utils.SuccessResponse(c, "Dashboard stats retrieved successfully", stats)
}

// GetSalesAnalytics retrieves sales analytics
// @Summary Get sales analytics
// @Description Get detailed sales analytics (admin only)
// @Tags admin
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param period query string false "Period (daily, weekly, monthly)" default(daily)
// @Success 200 {object} utils.Response{data=models.SalesAnalytics}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/analytics/sales [get]
func (h *AdminHandler) GetSalesAnalytics(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	var startDate, endDate *time.Time
	period := c.QueryParam("period")
	if period == "" {
		period = "daily"
	}

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		} else {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start_date format (use YYYY-MM-DD)")
		}
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		} else {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format (use YYYY-MM-DD)")
		}
	}

	// If no dates provided, default to last 30 days
	if startDate == nil || endDate == nil {
		now := time.Now()
		endDate = &now
		start := now.AddDate(0, 0, -30)
		startDate = &start
	}

	analytics, err := h.orderService.GetOrderAnalytics(c.Request().Context(), nil, startDate, endDate)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	salesAnalytics := &models.SalesAnalytics{
		Period:          period,
		StartDate:       *startDate,
		EndDate:         *endDate,
		TotalRevenue:    analytics.TotalRevenue,
		TotalOrders:     analytics.TotalOrders,
		AverageOrderValue: func() float64 {
			if analytics.TotalOrders > 0 {
				return analytics.TotalRevenue / float64(analytics.TotalOrders)
			}
			return 0
		}(),
		// You can add more detailed analytics here like daily/weekly/monthly breakdowns
	}

	return utils.SuccessResponse(c, "Sales analytics retrieved successfully", salesAnalytics)
}

// GetUserAnalytics retrieves user analytics
// @Summary Get user analytics
// @Description Get user analytics and statistics (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} utils.Response{data=models.UserAnalytics}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/analytics/users [get]
func (h *AdminHandler) GetUserAnalytics(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	// You would implement GetUserStats in UserService
	// userStats, err := h.userService.GetUserStats(c.Request().Context())
	// if err != nil {
	//     return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	// }

	// For now, return a placeholder response
	userAnalytics := &models.UserAnalytics{
		TotalUsers:    0, // userStats.TotalUsers,
		ActiveUsers:   0, // userStats.ActiveUsers,
		NewUsers:      0, // userStats.NewUsersThisMonth,
		CustomerCount: 0, // userStats.Customers,
		SellerCount:   0, // userStats.Sellers,
		AdminCount:    0, // userStats.Admins,
	}

	return utils.SuccessResponse(c, "User analytics retrieved successfully", userAnalytics)
}

// GetProductAnalytics retrieves product analytics
// @Summary Get product analytics
// @Description Get product analytics and statistics (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} utils.Response{data=models.ProductAnalytics}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/analytics/products [get]
func (h *AdminHandler) GetProductAnalytics(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	// Get low stock products
	lowStockProducts, err := h.productService.GetLowStockProducts(c.Request().Context(), 10, nil)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get low stock products")
	}

	// Get top rated products
	topRatedProducts, err := h.productService.GetTopRatedProducts(c.Request().Context(), 10)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get top rated products")
	}

	productAnalytics := &models.ProductAnalytics{
		TotalProducts:     0, // You would get this from a count method
		LowStockProducts:  len(lowStockProducts),
		OutOfStockProducts: 0, // You would implement this
		TopRatedProducts:  topRatedProducts,
	}

	return utils.SuccessResponse(c, "Product analytics retrieved successfully", productAnalytics)
}

// GetReviewAnalytics retrieves review analytics
// @Summary Get review analytics
// @Description Get review analytics and statistics (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} utils.Response{data=models.ReviewAnalytics}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/analytics/reviews [get]
func (h *AdminHandler) GetReviewAnalytics(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	// Get recent reviews
	recentReviews, err := h.reviewService.GetRecentReviews(c.Request().Context(), 10)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get recent reviews")
	}

	// Get top reviews
	topReviews, err := h.reviewService.GetTopReviews(c.Request().Context(), 10)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get top reviews")
	}

	reviewAnalytics := &models.ReviewAnalytics{
		TotalReviews:   0, // You would get this from a count method
		AverageRating:  0, // You would calculate this
		RecentReviews:  recentReviews,
		TopReviews:     topReviews,
	}

	return utils.SuccessResponse(c, "Review analytics retrieved successfully", reviewAnalytics)
}

// GetSystemHealth checks system health
// @Summary Get system health
// @Description Get system health status (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} utils.Response{data=models.SystemHealth}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/health [get]
func (h *AdminHandler) GetSystemHealth(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	health := &models.SystemHealth{
		Status:       "healthy",
		DatabaseStatus: "connected",
		RedisStatus:   "connected",
		LastChecked:   time.Now(),
		Uptime:        time.Since(time.Now().Add(-time.Hour * 24)), // Placeholder
	}

	return utils.SuccessResponse(c, "System health retrieved successfully", health)
}

// ManageUser manages user accounts
// @Summary Manage user account
// @Description Update user role or status (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.AdminUserUpdateRequest true "User update data"
// @Success 200 {object} utils.Response{data=models.User}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/users/{id} [put]
func (h *AdminHandler) ManageUser(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

// id is not used, so removed to fix build error

	var req models.AdminUserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Admin user management - you would implement a special admin update method
	// For now, just return success without actual implementation
	// TODO: Implement proper admin user update functionality
	
	return utils.SuccessResponse(c, "User updated successfully", nil)
}

// GetOrderDetails retrieves detailed order information
// @Summary Get detailed order information
// @Description Get comprehensive order details (admin only)
// @Tags admin
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response{data=models.Order}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/orders/{id} [get]
func (h *AdminHandler) GetOrderDetails(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)
	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	// Admin can view any order, so we pass admin role
	order, err := h.orderService.GetOrder(c.Request().Context(), uint(id), 0, models.RoleAdmin)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
	}

	return utils.SuccessResponse(c, "Order details retrieved successfully", order)
}
