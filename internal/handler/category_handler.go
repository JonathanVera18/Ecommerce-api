package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/service"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req models.CategoryCreateRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	category, err := h.categoryService.CreateCategory(c.Request().Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, "Category created successfully", category)
}

// GetCategory retrieves a category by ID
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
	}

	category, err := h.categoryService.GetCategory(c.Request().Context(), uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
	}

	return utils.SuccessResponse(c, "Category retrieved successfully", category)
}

// GetAllCategories retrieves all categories
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	categories, err := h.categoryService.GetAllCategories(c.Request().Context())
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Categories retrieved successfully", categories)
}

// GetCategoryBySlug retrieves a category by slug
func (h *CategoryHandler) GetCategoryBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Slug is required")
	}

	category, err := h.categoryService.GetCategoryBySlug(c.Request().Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
	}

	return utils.SuccessResponse(c, "Category retrieved successfully", category)
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
	}

	var req models.CategoryUpdateRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	category, err := h.categoryService.UpdateCategory(c.Request().Context(), uint(id), &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Category updated successfully", category)
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
	}

	err = h.categoryService.DeleteCategory(c.Request().Context(), uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Category deleted successfully", nil)
}

// GetCategoriesHierarchy retrieves categories in hierarchical structure
func (h *CategoryHandler) GetCategoriesHierarchy(c echo.Context) error {
	categories, err := h.categoryService.GetCategoriesHierarchy(c.Request().Context())
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Categories hierarchy retrieved successfully", categories)
}

// GetCategoryChildren retrieves child categories of a parent category
func (h *CategoryHandler) GetCategoryChildren(c echo.Context) error {
	parentID, err := strconv.ParseUint(c.Param("parentId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid parent category ID")
	}

	categories, err := h.categoryService.GetCategoryChildren(c.Request().Context(), uint(parentID))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, "Child categories retrieved successfully", categories)
}
