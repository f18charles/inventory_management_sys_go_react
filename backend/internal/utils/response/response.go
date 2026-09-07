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

func NewPaginationMeta(page, pageSize int, totalItems int64) PaginationMeta {
	if pageSize <= 0 {
		pageSize = 20
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	return PaginationMeta{Page: page, PageSize: pageSize, TotalItems: totalItems, TotalPages: totalPages}
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

func Paginated(c *gin.Context, status int, data any, meta PaginationMeta) {
	c.JSON(status, gin.H{"data": data, "meta": meta})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": ErrorBody{Code: code, Message: message}})
}
