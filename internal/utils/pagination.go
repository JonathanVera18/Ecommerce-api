package utils

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

// PaginationParams extracts pagination parameters from query string
func PaginationParams(c echo.Context) (page, limit int) {
	page = 1
	limit = 20

	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit
}

// BuildPaginationMeta creates pagination metadata
func BuildPaginationMeta(page, limit int, total int64) *models.PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	
	return &models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// GetOffset calculates the database offset for pagination
func GetOffset(page, limit int) int {
	return (page - 1) * limit
}
