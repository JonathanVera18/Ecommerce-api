package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// AddToCart adds a product to user's cart
func (h *CartHandler) AddToCart(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req models.CartAddRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	cart, err := h.cartService.AddToCart(c.Request().Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, "Product added to cart successfully", cart)
}

// UpdateCartItem updates quantity of a product in user's cart
func (h *CartHandler) UpdateCartItem(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req struct {
		Quantity int `json:"quantity" validate:"required,min=1"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	cart, err := h.cartService.UpdateCartItem(c.Request().Context(), userID, uint(productID), req.Quantity)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Cart item updated successfully", cart)
}

// RemoveFromCart removes a product from user's cart
func (h *CartHandler) RemoveFromCart(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	err = h.cartService.RemoveFromCart(c.Request().Context(), userID, uint(productID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product removed from cart successfully", nil)
}

// GetUserCart retrieves user's cart
func (h *CartHandler) GetUserCart(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	cart, err := h.cartService.GetUserCart(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

// GetCartTotal retrieves user's cart total
func (h *CartHandler) GetCartTotal(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	total, err := h.cartService.GetCartTotal(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Cart total retrieved successfully", map[string]float64{"total": total})
}

// ClearCart clears user's entire cart
func (h *CartHandler) ClearCart(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	err := h.cartService.ClearCart(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Cart cleared successfully", nil)
}

// GetCartItemCount retrieves user's cart item count
func (h *CartHandler) GetCartItemCount(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	count, err := h.cartService.GetCartItemCount(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Cart item count retrieved successfully", map[string]int{"count": count})
}
