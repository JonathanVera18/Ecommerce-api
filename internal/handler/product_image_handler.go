package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type ProductImageHandler struct {
	productImageService service.ProductImageService
}

func NewProductImageHandler(productImageService service.ProductImageService) *ProductImageHandler {
	return &ProductImageHandler{
		productImageService: productImageService,
	}
}

// AddProductImage adds a new image to a product
// @Summary Add product image
// @Description Add a new image to a product
// @Tags product-images
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param image body models.ProductImageRequest true "Image data"
// @Success 201 {object} utils.Response{data=models.ProductImage}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/images [post]
func (h *ProductImageHandler) AddProductImage(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req models.ProductImageRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	image, err := h.productImageService.AddProductImage(c.Request().Context(), uint(productID), &req)
	if err != nil {
		if err.Error() == "product not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, "Product image added successfully", image)
}

// GetProductImages gets all images for a product
// @Summary Get product images
// @Description Get all images for a product
// @Tags product-images
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} utils.Response{data=[]models.ProductImage}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/{product_id}/images [get]
func (h *ProductImageHandler) GetProductImages(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	images, err := h.productImageService.GetProductImages(c.Request().Context(), uint(productID))
	if err != nil {
		if err.Error() == "product not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product images retrieved successfully", images)
}

// GetProductImage gets a specific product image
// @Summary Get product image
// @Description Get a specific product image by ID
// @Tags product-images
// @Produce json
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Success 200 {object} utils.Response{data=models.ProductImage}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/{product_id}/images/{image_id} [get]
func (h *ProductImageHandler) GetProductImage(c echo.Context) error {
	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid image ID")
	}

	image, err := h.productImageService.GetProductImage(c.Request().Context(), uint(imageID))
	if err != nil {
		if err.Error() == "product image not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product image retrieved successfully", image)
}

// UpdateProductImage updates a product image
// @Summary Update product image
// @Description Update a product image
// @Tags product-images
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Param image body models.ProductImageRequest true "Image data"
// @Success 200 {object} utils.Response{data=models.ProductImage}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/images/{image_id} [put]
func (h *ProductImageHandler) UpdateProductImage(c echo.Context) error {
	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid image ID")
	}

	var req models.ProductImageRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	image, err := h.productImageService.UpdateProductImage(c.Request().Context(), uint(imageID), &req)
	if err != nil {
		if err.Error() == "product image not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product image updated successfully", image)
}

// DeleteProductImage deletes a product image
// @Summary Delete product image
// @Description Delete a product image
// @Tags product-images
// @Produce json
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/images/{image_id} [delete]
func (h *ProductImageHandler) DeleteProductImage(c echo.Context) error {
	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid image ID")
	}

	err = h.productImageService.DeleteProductImage(c.Request().Context(), uint(imageID))
	if err != nil {
		if err.Error() == "product image not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product image deleted successfully", nil)
}

// SetPrimaryImage sets an image as the primary image for a product
// @Summary Set primary image
// @Description Set an image as the primary image for a product
// @Tags product-images
// @Produce json
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/images/{image_id}/primary [put]
func (h *ProductImageHandler) SetPrimaryImage(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid image ID")
	}

	err = h.productImageService.SetPrimaryImage(c.Request().Context(), uint(productID), uint(imageID))
	if err != nil {
		if err.Error() == "product not found" || err.Error() == "image not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		if err.Error() == "image does not belong to the specified product" {
			return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Primary image set successfully", nil)
}

// GetPrimaryImage gets the primary image for a product
// @Summary Get primary image
// @Description Get the primary image for a product
// @Tags product-images
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} utils.Response{data=models.ProductImage}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /products/{product_id}/images/primary [get]
func (h *ProductImageHandler) GetPrimaryImage(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	image, err := h.productImageService.GetPrimaryImage(c.Request().Context(), uint(productID))
	if err != nil {
		if err.Error() == "product not found" || err.Error() == "primary image not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Primary image retrieved successfully", image)
}

// UpdateImageOrder updates the sort order of an image
// @Summary Update image order
// @Description Update the sort order of an image
// @Tags product-images
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Param order body map[string]int true "Sort order"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/images/{image_id}/order [put]
func (h *ProductImageHandler) UpdateImageOrder(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	imageID, err := strconv.ParseUint(c.Param("image_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid image ID")
	}

	var req map[string]int
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	sortOrder, exists := req["sort_order"]
	if !exists {
		return utils.ErrorResponse(c, http.StatusBadRequest, "sort_order is required")
	}

	err = h.productImageService.UpdateImageOrder(c.Request().Context(), uint(productID), uint(imageID), sortOrder)
	if err != nil {
		if err.Error() == "product not found" || err.Error() == "image not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		if err.Error() == "image does not belong to the specified product" {
			return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Image order updated successfully", nil)
}

// BulkAddImages adds multiple images to a product
// @Summary Bulk add images
// @Description Add multiple images to a product
// @Tags product-images
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param images body []models.ProductImageRequest true "Images data"
// @Success 201 {object} utils.Response{data=[]models.ProductImage}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/images/bulk [post]
func (h *ProductImageHandler) BulkAddImages(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req []models.ProductImageRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if len(req) == 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "At least one image is required")
	}

	// Validate each image request
	for i, imageReq := range req {
		if err := utils.ValidateStruct(&imageReq); err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid image at index %d: %s", i, err.Error()))
		}
	}

	images, err := h.productImageService.BulkAddImages(c.Request().Context(), uint(productID), req)
	if err != nil {
		if err.Error() == "product not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, "Images added successfully", images)
}

// ReplaceProductImages replaces all images for a product
// @Summary Replace product images
// @Description Replace all images for a product
// @Tags product-images
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param images body []models.ProductImageRequest true "Images data"
// @Success 200 {object} utils.Response{data=[]models.ProductImage}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /products/{product_id}/images/replace [put]
func (h *ProductImageHandler) ReplaceProductImages(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req []models.ProductImageRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Validate each image request
	for i, imageReq := range req {
		if err := utils.ValidateStruct(&imageReq); err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid image at index %d: %s", i, err.Error()))
		}
	}

	images, err := h.productImageService.ReplaceProductImages(c.Request().Context(), uint(productID), req)
	if err != nil {
		if err.Error() == "product not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Product images replaced successfully", images)
}
