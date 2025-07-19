package service

import (
	"context"
	"errors"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"gorm.io/gorm"
)

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}



func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) AddToCart(ctx context.Context, userID uint, req *models.CartAddRequest) (*models.CartResponse, error) {
	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Check if product is in stock
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.GetItemByProduct(ctx, cart.ID, req.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingItem != nil {
		// Update quantity
		existingItem.Quantity += req.Quantity
		if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		// Add new item
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := s.cartRepo.AddItem(ctx, cartItem); err != nil {
			return nil, err
		}
	}

	// Return updated cart
	return s.GetCart(ctx, userID)
}

func (s *cartService) GetCart(ctx context.Context, userID uint) (*models.CartResponse, error) {
	cart, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty cart
			return &models.CartResponse{
				CustomerID:  userID,
				Items:       []models.CartItemResponse{},
				TotalAmount: 0,
				ItemCount:   0,
			}, nil
		}
		return nil, err
	}

	resp := cart.ToResponse()
	return &resp, nil
}

func (s *cartService) ClearCart(ctx context.Context, userID uint) error {
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	return s.cartRepo.ClearCart(ctx, cart.ID)
}

func (s *cartService) GetUserCart(ctx context.Context, userID uint) ([]*models.CartResponse, error) {
	cartWithItems, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []*models.CartResponse
	for _, item := range cartWithItems.CartItems {
		// Get product details
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			continue
		}

		itemResponse := models.CartItemResponse{
			ID:        item.ID,
			CartID:    item.CartID,
			ProductID: item.ProductID,
			Product:   product.ToResponse(),
			Quantity:  item.Quantity,
			Subtotal:  product.Price * float64(item.Quantity),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		cartResponse := &models.CartResponse{
			ID:          cartWithItems.ID,
			CustomerID:  cartWithItems.CustomerID,
			Items:       []models.CartItemResponse{itemResponse},
			TotalAmount: itemResponse.Subtotal,
			ItemCount:   item.Quantity,
			CreatedAt:   cartWithItems.CreatedAt,
			UpdatedAt:   cartWithItems.UpdatedAt,
		}

		responses = append(responses, cartResponse)
	}

	return responses, nil
}

func (s *cartService) GetCartTotal(ctx context.Context, userID uint) (float64, error) {
	cartWithItems, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range cartWithItems.CartItems {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			continue
		}
		total += product.Price * float64(item.Quantity)
	}

	return total, nil
}

func (s *cartService) GetCartItemCount(ctx context.Context, userID uint) (int, error) {
	cartWithItems, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return 0, err
	}

	var count int
	for _, item := range cartWithItems.CartItems {
		count += item.Quantity
	}

	return count, nil
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID uint, productID uint, quantity int) (*models.CartResponse, error) {
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	existingItem, err := s.cartRepo.GetItemByProduct(ctx, cart.ID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found in cart")
		}
		return nil, err
	}

	// Check product stock
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product.Stock < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Update quantity
	existingItem.Quantity = quantity
	if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
		return nil, err
	}

	return &models.CartResponse{
		ID:          cart.ID,
		CustomerID:  cart.CustomerID,
		Items:       []models.CartItemResponse{},
		TotalAmount: product.Price * float64(existingItem.Quantity),
		ItemCount:   existingItem.Quantity,
		CreatedAt:   existingItem.CreatedAt,
		UpdatedAt:   existingItem.UpdatedAt,
	}, nil
}

func (s *cartService) RemoveFromCart(ctx context.Context, userID uint, productID uint) error {
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	return s.cartRepo.RemoveItem(ctx, cart.ID, productID)
}
