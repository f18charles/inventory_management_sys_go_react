// internal/utils/response/response.go  (+ http.go)
//
// The ONLY place the API's JSON envelope is constructed, plus the shared
// request-binding and domain-error → HTTP mapping used by every domain handler.
// Handlers call these instead of building gin.H{...} inline, so the shape can
// never drift between endpoints.

package response

import (
	"errors"
	"math"
	"net/http"

	"i_m_s/internal/models"

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

// NewPaginationMeta is pure arithmetic over a COUNT result the caller already
// has — nothing here talks to the database.
func NewPaginationMeta(page, pageSize int, totalItems int64) PaginationMeta {
	if pageSize <= 0 {
		pageSize = 20
	}
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

// RespondError is the single mapping from domain sentinel errors to HTTP status
// + envelope. Handlers call this instead of hand-rolling a switch each time.
// The default case never leaks raw SQL/GORM/stack detail to the caller.
func RespondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		Error(c, http.StatusNotFound, "NOT_FOUND", "resource not found")
	case errors.Is(err, models.ErrInsufficientInventory):
		Error(c, http.StatusConflict, "INSUFFICIENT_INVENTORY", err.Error())
	case errors.Is(err, models.ErrUnauthorized):
		Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	case errors.Is(err, models.ErrInvalidCredentials):
		Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password")
	case errors.Is(err, models.ErrAccountInactive):
		Error(c, http.StatusForbidden, "ACCOUNT_INACTIVE", "account is inactive")
	case errors.Is(err, models.ErrForbidden):
		Error(c, http.StatusForbidden, "FORBIDDEN", "forbidden")
	case errors.Is(err, models.ErrBadRequest):
		Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	case errors.Is(err, models.ErrConflict):
		Error(c, http.StatusConflict, "CONFLICT", err.Error())
	case errors.Is(err, models.ErrValidationError):
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
	default:
		Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
	}
}

// BindJSON binds and validates a JSON body, writing the standard error envelope
// and returning false when binding fails.
func BindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return false
	}
	return true
}
