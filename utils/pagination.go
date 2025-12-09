package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination represents pagination parameters
type Pagination struct {
	Page  int
	Limit int
}

// GetPagination extracts pagination parameters from query string
func GetPagination(c *gin.Context) Pagination {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // max limit
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset calculates the offset for database query
func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// CalculateTotalPages calculates total pages based on total items and limit
func CalculateTotalPages(totalItems int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	pages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		pages++
	}
	return pages
}
