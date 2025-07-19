package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"github.com/JonathanVera18/ecommerce-api/pkg/payment"
)

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	paymentSvc  payment.Service
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	paymentSvc payment.Service,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		paymentSvc:  paymentSvc,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *models.CreateOrderRequest, userID uint) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	// Validate and calculate order items
	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}

		if !product.IsActive {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (available: %d, requested: %d)",
				product.Name, product.Stock, item.Quantity)
		}

		itemTotal := product.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:          item.ProductID,
			Quantity:           item.Quantity,
			UnitPrice:          product.Price,
			TotalPrice:         itemTotal,
			ProductName:        product.Name,
			ProductSKU:         product.SKU,
			ProductDescription: &product.Description,
			ProductImage:       nil, // Will be set from product images if available
		})
	}

	// Create order
	order := &models.Order{
		CustomerID:         userID,
		Status:             models.OrderStatusPending,
		TotalAmount:        totalAmount,
		SubtotalAmount:     totalAmount,
		PaymentMethod:      req.PaymentMethod,
		ShippingFirstName:  "Customer", // These should come from user profile or request
		ShippingLastName:   "User",
		ShippingEmail:      "customer@example.com",
		ShippingStreet:     req.ShippingAddress,
		ShippingCity:       "City",
		ShippingState:      "State", 
		ShippingCountry:    "Country",
		ShippingPostalCode: "12345",
		OrderItems:         orderItems,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Update product stock
	for _, item := range req.Items {
		product, _ := s.productRepo.GetByID(ctx, item.ProductID)
		newStock := product.Stock - item.Quantity
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			// Log error but don't fail the order creation
			// In production, you might want to implement a rollback mechanism
			fmt.Printf("Warning: failed to update stock for product %d: %v\n", item.ProductID, err)
		}
	}

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id uint, userID uint, userRole models.UserRole) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Check authorization
	if userRole != models.RoleAdmin && order.CustomerID != userID {
		// Check if user is seller of any item in the order
		if userRole == models.RoleSeller {
			hasSellerItem := false
			for _, item := range order.OrderItems {
				product, err := s.productRepo.GetByID(ctx, item.ProductID)
				if err == nil && product.SellerID == userID {
					hasSellerItem = true
					break
				}
			}
			if !hasSellerItem {
				return nil, errors.New("unauthorized to view this order")
			}
		} else {
			return nil, errors.New("unauthorized to view this order")
		}
	}

	return order, nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID uint, limit, offset int) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	return orders, nil
}

func (s *orderService) GetAllOrders(ctx context.Context, limit, offset int) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}

	return orders, nil
}

func (s *orderService) GetOrdersByStatus(ctx context.Context, status models.OrderStatus, limit, offset int) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}

	return orders, nil
}

func (s *orderService) GetSellerOrders(ctx context.Context, sellerID uint, limit, offset int) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetOrdersBySellerID(ctx, sellerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller orders: %w", err)
	}

	return orders, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id uint, status models.OrderStatus, userID uint, userRole models.UserRole) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Check authorization for status updates
	if userRole != models.RoleAdmin {
		if userRole == models.RoleSeller {
			// Sellers can only update orders containing their products
			hasSellerItem := false
			for _, item := range order.OrderItems {
				product, err := s.productRepo.GetByID(ctx, item.ProductID)
				if err == nil && product.SellerID == userID {
					hasSellerItem = true
					break
				}
			}
			if !hasSellerItem {
				return errors.New("unauthorized to update this order")
			}
		} else {
			return errors.New("unauthorized to update order status")
		}
	}

	// Validate status transition
	if !isValidStatusTransition(order.Status, status) {
		return fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *orderService) ProcessPayment(ctx context.Context, orderID uint, paymentReq *models.PaymentRequest) (*models.PaymentResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != models.OrderStatusPending {
		return nil, errors.New("order is not in pending status")
	}

	// Process payment using payment service
	paymentIntentID, err := s.paymentSvc.CreatePaymentIntent(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	// Confirm payment
	err = s.paymentSvc.ConfirmPayment(paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("payment confirmation failed: %w", err)
	}

	// Update order status to confirmed
	if err := s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusConfirmed); err != nil {
		return nil, fmt.Errorf("failed to update order status after payment: %w", err)
	}

	return &models.PaymentResponse{
		TransactionID: paymentIntentID,
		Status:        "confirmed",
		Amount:        order.TotalAmount,
	}, nil
}

func (s *orderService) CancelOrder(ctx context.Context, id uint, userID uint, userRole models.UserRole) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Check authorization
	if userRole != models.RoleAdmin && order.CustomerID != userID {
		return errors.New("unauthorized to cancel this order")
	}

	// Can only cancel pending or confirmed orders
	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusConfirmed {
		return errors.New("order cannot be cancelled in its current status")
	}

	// Restore product stock
	for _, item := range order.OrderItems {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err == nil {
			newStock := product.Stock + item.Quantity
			s.productRepo.UpdateStock(ctx, item.ProductID, newStock)
		}
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, models.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

func (s *orderService) GetOrderAnalytics(ctx context.Context, sellerID *uint, startDate, endDate *time.Time) (*models.OrderAnalytics, error) {
	var totalRevenue float64
	var err error

	if sellerID != nil {
		totalRevenue, err = s.orderRepo.GetRevenueBySellerID(ctx, *sellerID, startDate, endDate)
	} else {
		totalRevenue, err = s.orderRepo.GetTotalRevenue(ctx, startDate, endDate)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get revenue: %w", err)
	}

	// Get total orders count
	totalOrders, err := s.orderRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order count: %w", err)
	}

	// Get orders by status
	pendingOrders, _ := s.orderRepo.CountByStatus(ctx, models.OrderStatusPending)
	confirmedOrders, _ := s.orderRepo.CountByStatus(ctx, models.OrderStatusConfirmed)
	shippedOrders, _ := s.orderRepo.CountByStatus(ctx, models.OrderStatusShipped)
	deliveredOrders, _ := s.orderRepo.CountByStatus(ctx, models.OrderStatusDelivered)
	cancelledOrders, _ := s.orderRepo.CountByStatus(ctx, models.OrderStatusCancelled)

	return &models.OrderAnalytics{
		TotalRevenue:     totalRevenue,
		TotalOrders:      totalOrders,
		PendingOrders:    pendingOrders,
		ConfirmedOrders:  confirmedOrders,
		ShippedOrders:    shippedOrders,
		DeliveredOrders:  deliveredOrders,
		CancelledOrders:  cancelledOrders,
	}, nil
}

func isValidStatusTransition(from, to models.OrderStatus) bool {
	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPending:   {models.OrderStatusConfirmed, models.OrderStatusCancelled},
		models.OrderStatusConfirmed: {models.OrderStatusProcessing, models.OrderStatusCancelled},
		models.OrderStatusProcessing: {models.OrderStatusShipped, models.OrderStatusCancelled},
		models.OrderStatusShipped:   {models.OrderStatusDelivered},
		models.OrderStatusDelivered: {}, // Final state
		models.OrderStatusCancelled: {}, // Final state
	}

	validStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, validState := range validStates {
		if validState == to {
			return true
		}
	}

	return false
}
