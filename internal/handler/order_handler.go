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

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder creates a new order
// @Summary Create a new order
// @Description Create a new order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order data"
// @Success 201 {object} utils.Response{data=models.Order}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req models.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.GetValidationErrors(err))
	}

	order, err := h.orderService.CreateOrder(c.Request().Context(), &req, userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Order created successfully", order)
}

// GetOrder retrieves an order by ID
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response{data=models.Order}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	order, err := h.orderService.GetOrder(c.Request().Context(), uint(id), userID, userRole)
	if err != nil {
		if err.Error() == "unauthorized to view this order" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
	}

	return utils.SuccessResponse(c, "Order retrieved successfully", order)
}

// GetUserOrders retrieves orders for the current user
// @Summary Get user orders
// @Description Get orders for the authenticated user
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Order}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders/my [get]
func (h *OrderHandler) GetUserOrders(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	orders, err := h.orderService.GetUserOrders(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Orders retrieved successfully", orders)
}

// GetAllOrders retrieves all orders (admin only)
// @Summary Get all orders
// @Description Get all orders (admin only)
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Order}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /admin/orders [get]
func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	orders, err := h.orderService.GetAllOrders(c.Request().Context(), limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Orders retrieved successfully", orders)
}

// GetOrdersByStatus retrieves orders by status
// @Summary Get orders by status
// @Description Get orders filtered by status (admin/seller)
// @Tags orders
// @Produce json
// @Param status path string true "Order status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Order}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders/status/{status} [get]
func (h *OrderHandler) GetOrdersByStatus(c echo.Context) error {
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleAdmin && userRole != models.RoleSeller {
		return utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
	}

	statusStr := c.Param("status")
	status := models.OrderStatus(statusStr)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	orders, err := h.orderService.GetOrdersByStatus(c.Request().Context(), status, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Orders retrieved successfully", orders)
}

// GetSellerOrders retrieves orders for a seller
// @Summary Get seller orders
// @Description Get orders containing seller's products
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Order}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /seller/orders [get]
func (h *OrderHandler) GetSellerOrders(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleSeller && userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Seller access required")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	orders, err := h.orderService.GetSellerOrders(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Seller orders retrieved successfully", orders)
}

// UpdateOrderStatus updates the status of an order
// @Summary Update order status
// @Description Update order status (admin/seller)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "Status update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	var req models.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	err = h.orderService.UpdateOrderStatus(c.Request().Context(), uint(id), req.Status, userID, userRole)
	if err != nil {
		if err.Error() == "unauthorized to update this order" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Order status updated successfully", nil)
}

// ProcessPayment processes payment for an order
// @Summary Process payment
// @Description Process payment for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param payment body models.PaymentRequest true "Payment data"
// @Success 200 {object} utils.Response{data=models.PaymentResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders/{id}/payment [post]
func (h *OrderHandler) ProcessPayment(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	var req models.PaymentRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.GetValidationErrors(err))
	}

	paymentResponse, err := h.orderService.ProcessPayment(c.Request().Context(), uint(id), &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Payment processed successfully", paymentResponse)
}

// CancelOrder cancels an order
// @Summary Cancel order
// @Description Cancel an order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	err = h.orderService.CancelOrder(c.Request().Context(), uint(id), userID, userRole)
	if err != nil {
		if err.Error() == "unauthorized to cancel this order" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Order cancelled successfully", nil)
}

// GetOrderAnalytics retrieves order analytics
// @Summary Get order analytics
// @Description Get order analytics (admin/seller)
// @Tags orders
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} utils.Response{data=models.OrderAnalytics}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /orders/analytics [get]
func (h *OrderHandler) GetOrderAnalytics(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleAdmin && userRole != models.RoleSeller {
		return utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
	}

	var startDate, endDate *time.Time

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	var sellerID *uint
	if userRole == models.RoleSeller {
		sellerID = &userID
	}

	analytics, err := h.orderService.GetOrderAnalytics(c.Request().Context(), sellerID, startDate, endDate)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Order analytics retrieved successfully", analytics)
}
