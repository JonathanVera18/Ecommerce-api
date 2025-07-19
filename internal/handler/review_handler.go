package handler

import (
	"net/http"
	"strconv"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// CreateReview creates a new review
// @Summary Create a new review
// @Description Create a new product review
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.CreateReviewRequest true "Review data"
// @Success 201 {object} utils.Response{data=models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req models.CreateReviewRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		validationErrors := utils.GetValidationErrors(err)
		return utils.ValidationError(c, validationErrors)
	}

	review, err := h.reviewService.CreateReview(c.Request().Context(), &req, userID)
	if err != nil {
		if err.Error() == "you can only review products you have purchased and received" ||
			err.Error() == "you have already reviewed this product" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Review created successfully", review)
}

// GetReview retrieves a review by ID
// @Summary Get review by ID
// @Description Get review details by ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} utils.Response{data=models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
	}

	review, err := h.reviewService.GetReview(c.Request().Context(), uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Review not found")
	}

	return utils.SuccessResponse(c, "Review retrieved successfully", review)
}

// GetProductReviews retrieves reviews for a product
// @Summary Get product reviews
// @Description Get reviews for a specific product
// @Tags reviews
// @Produce json
// @Param product_id path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/{product_id}/reviews [get]
func (h *ReviewHandler) GetProductReviews(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
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

	reviews, err := h.reviewService.GetProductReviews(c.Request().Context(), uint(productID), limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product reviews retrieved successfully", reviews)
}

// GetUserReviews retrieves reviews by a user
// @Summary Get user reviews
// @Description Get reviews written by the authenticated user
// @Tags reviews
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /reviews/my [get]
func (h *ReviewHandler) GetUserReviews(c echo.Context) error {
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

	reviews, err := h.reviewService.GetUserReviews(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "User reviews retrieved successfully", reviews)
}

// UpdateReview updates an existing review
// @Summary Update a review
// @Description Update an existing review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body models.UpdateReviewRequest true "Updated review data"
// @Success 200 {object} utils.Response{data=models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
	}

	var req models.UpdateReviewRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	review, err := h.reviewService.UpdateReview(c.Request().Context(), uint(id), &req, userID)
	if err != nil {
		if err.Error() == "unauthorized to update this review" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Review updated successfully", review)
}

// DeleteReview deletes a review
// @Summary Delete a review
// @Description Delete a review (user who created it or admin)
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
	}

	err = h.reviewService.DeleteReview(c.Request().Context(), uint(id), userID, userRole)
	if err != nil {
		if err.Error() == "unauthorized to delete this review" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Review deleted successfully", nil)
}

// GetReviewsByRating retrieves reviews by rating
// @Summary Get reviews by rating
// @Description Get reviews filtered by rating
// @Tags reviews
// @Produce json
// @Param rating path int true "Rating (1-5)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /reviews/rating/{rating} [get]
func (h *ReviewHandler) GetReviewsByRating(c echo.Context) error {
	rating, err := strconv.Atoi(c.Param("rating"))
	if err != nil || rating < 1 || rating > 5 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid rating (must be 1-5)")
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

	reviews, err := h.reviewService.GetReviewsByRating(c.Request().Context(), rating, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Reviews by rating retrieved successfully", reviews)
}

// GetTopReviews retrieves top-rated reviews
// @Summary Get top reviews
// @Description Get highest rated reviews
// @Tags reviews
// @Produce json
// @Param limit query int false "Number of reviews to return" default(10)
// @Success 200 {object} utils.Response{data=[]models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /reviews/top [get]
func (h *ReviewHandler) GetTopReviews(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	reviews, err := h.reviewService.GetTopReviews(c.Request().Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Top reviews retrieved successfully", reviews)
}

// GetRecentReviews retrieves recent reviews
// @Summary Get recent reviews
// @Description Get most recent reviews
// @Tags reviews
// @Produce json
// @Param limit query int false "Number of reviews to return" default(10)
// @Success 200 {object} utils.Response{data=[]models.Review}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /reviews/recent [get]
func (h *ReviewHandler) GetRecentReviews(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	reviews, err := h.reviewService.GetRecentReviews(c.Request().Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Recent reviews retrieved successfully", reviews)
}

// GetProductReviewStats retrieves review statistics for a product
// @Summary Get product review stats
// @Description Get review statistics for a specific product
// @Tags reviews
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} utils.Response{data=models.ReviewStats}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/{product_id}/reviews/stats [get]
func (h *ReviewHandler) GetProductReviewStats(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	stats, err := h.reviewService.GetProductReviewStats(c.Request().Context(), uint(productID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product review stats retrieved successfully", stats)
}

// CanUserReview checks if user can review a product
// @Summary Check if user can review
// @Description Check if authenticated user can review a specific product
// @Tags reviews
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} utils.Response{data=bool}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/can-review [get]
func (h *ReviewHandler) CanUserReview(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	canReview, err := h.reviewService.CanUserReview(c.Request().Context(), userID, uint(productID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Review eligibility checked successfully", map[string]bool{
		"can_review": canReview,
	})
}
