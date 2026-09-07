// internal/utils/response/response.go
//
// The ONLY place that constructs the API's JSON envelope. Handlers call
// these instead of building gin.H{...} inline, so the shape can never
// drift between endpoints.

package response

import (
	"math"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginationMeta computes total pages from a total row count — callers
// pass what they already have (page, pageSize, totalItems from a COUNT
// query); nothing here talks to the database.
func NewPaginationMeta(page, pageSize int, totalItems int64) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	return PaginationMeta{Page: page, PageSize: pageSize, TotalItems: totalItems, TotalPages: totalPages}
}

// Success: { "data": ... }
func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

// Paginated: { "data": [...], "meta": { "page": 1, "page_size": 20, ... } }
func Paginated(c *gin.Context, status int, data any, meta PaginationMeta) {
	c.JSON(status, gin.H{"data": data, "meta": meta})
}

// Error: { "error": { "code": "...", "message": "..." } }
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": ErrorBody{Code: code, Message: message}})
}
