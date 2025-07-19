package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type WishlistHandler struct {
	wishlistService service.WishlistService
}

func NewWishlistHandler(wishlistService service.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}

// AddToWishlist adds a product to user's wishlist
func (h *WishlistHandler) AddToWishlist(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req models.WishlistAddRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	wishlist, err := h.wishlistService.AddToWishlist(c.Request().Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, "Product added to wishlist successfully", wishlist)
}

// RemoveFromWishlist removes a product from user's wishlist
func (h *WishlistHandler) RemoveFromWishlist(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	err = h.wishlistService.RemoveFromWishlist(c.Request().Context(), userID, uint(productID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product removed from wishlist successfully", nil)
}

// GetUserWishlist retrieves user's wishlist
func (h *WishlistHandler) GetUserWishlist(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	wishlist, err := h.wishlistService.GetUserWishlist(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Wishlist retrieved successfully", wishlist)
}

// IsProductInWishlist checks if a product is in user's wishlist
func (h *WishlistHandler) IsProductInWishlist(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	isInWishlist, err := h.wishlistService.IsProductInWishlist(c.Request().Context(), userID, uint(productID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product wishlist status retrieved successfully", map[string]bool{"is_in_wishlist": isInWishlist})
}

// ClearWishlist clears user's entire wishlist
func (h *WishlistHandler) ClearWishlist(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	err := h.wishlistService.ClearWishlist(c.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Wishlist cleared successfully", nil)
}
