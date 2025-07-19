package handler

import (
	"net/http"
	"strconv"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Create a new product (seller only)
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Product data"
// @Success 201 {object} utils.Response{data=models.Product}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleSeller && userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Only sellers can create products")
	}

	var req models.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		validationErrors := utils.GetValidationErrors(err)
		return utils.ValidationError(c, validationErrors)
	}

	product, err := h.productService.CreateProduct(c.Request().Context(), &req, userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, "Product created successfully", product)
}

// GetProduct retrieves a product by ID
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response{data=models.Product}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.productService.GetProduct(c.Request().Context(), uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
	}

	return utils.SuccessResponse(c, "Product retrieved successfully", product)
}

// GetProducts retrieves products with filtering and pagination
// @Summary Get products
// @Description Get products with optional filtering
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query string false "Filter by category"
// @Param seller_id query int false "Filter by seller ID"
// @Param search query string false "Search in product name and description"
// @Success 200 {object} utils.Response{data=models.ProductListResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	req := &models.GetProductsRequest{
		Page:     page,
		Limit:    limit,
		Offset:   offset,
		Category: c.QueryParam("category"),
		Search:   c.QueryParam("search"),
	}

	if sellerIDStr := c.QueryParam("seller_id"); sellerIDStr != "" {
		if sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32); err == nil {
			sellerIDUint := uint(sellerID)
			req.SellerID = &sellerIDUint
		}
	}

	products, err := h.productService.GetProducts(c.Request().Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Products retrieved successfully", products)
}

// UpdateProduct updates an existing product
// @Summary Update a product
// @Description Update product details (seller/admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.UpdateProductRequest true "Updated product data"
// @Success 200 {object} utils.Response{data=models.Product}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleSeller && userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Only sellers can update products")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req models.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	product, err := h.productService.UpdateProduct(c.Request().Context(), uint(id), &req, userID)
	if err != nil {
		if err.Error() == "unauthorized to update this product" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product updated successfully", product)
}

// DeleteProduct deletes a product
// @Summary Delete a product
// @Description Delete a product (seller/admin only)
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleSeller && userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Only sellers can delete products")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	err = h.productService.DeleteProduct(c.Request().Context(), uint(id), userID)
	if err != nil {
		if err.Error() == "unauthorized to delete this product" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product deleted successfully", nil)
}

// UpdateStock updates product stock
// @Summary Update product stock
// @Description Update product stock quantity (seller/admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param stock body models.UpdateStockRequest true "Stock data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{id}/stock [put]
func (h *ProductHandler) UpdateStock(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleSeller && userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Only sellers can update stock")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req models.UpdateStockRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	err = h.productService.UpdateStock(c.Request().Context(), uint(id), req.Stock, userID)
	if err != nil {
		if err.Error() == "unauthorized to update this product's stock" {
			return utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Stock updated successfully", nil)
}

// GetLowStockProducts gets products with low stock
// @Summary Get low stock products
// @Description Get products with stock below threshold (seller/admin only)
// @Tags products
// @Produce json
// @Param threshold query int false "Stock threshold" default(10)
// @Success 200 {object} utils.Response{data=[]models.Product}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/low-stock [get]
func (h *ProductHandler) GetLowStockProducts(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	userRole := c.Get("user_role").(models.UserRole)

	if userRole != models.RoleSeller && userRole != models.RoleAdmin {
		return utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
	}

	threshold, _ := strconv.Atoi(c.QueryParam("threshold"))
	if threshold <= 0 {
		threshold = 10
	}

	var sellerID *uint
	if userRole == models.RoleSeller {
		sellerID = &userID
	}

	products, err := h.productService.GetLowStockProducts(c.Request().Context(), threshold, sellerID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Low stock products retrieved successfully", products)
}

// GetTopRatedProducts gets top rated products
// @Summary Get top rated products
// @Description Get products with highest ratings
// @Tags products
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} utils.Response{data=[]models.Product}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/top-rated [get]
func (h *ProductHandler) GetTopRatedProducts(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	products, err := h.productService.GetTopRatedProducts(c.Request().Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Top rated products retrieved successfully", products)
}

// SearchProducts searches for products
// @Summary Search products
// @Description Search products by name and description
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Product}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Search query is required")
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

	products, err := h.productService.SearchProducts(c.Request().Context(), query, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Search results retrieved successfully", products)
}

// GetProductsByCategory gets products by category
// @Summary Get products by category
// @Description Get products filtered by category
// @Tags products
// @Produce json
// @Param category path string true "Product category"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.Product}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/category/{category} [get]
func (h *ProductHandler) GetProductsByCategory(c echo.Context) error {
	category := c.Param("category")
	if category == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Category is required")
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

	products, err := h.productService.GetProductsByCategory(c.Request().Context(), category, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Products by category retrieved successfully", products)
}
