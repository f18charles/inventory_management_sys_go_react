package response

import (
	"errors"
	"net/http"

	"i_m_s/internal/models"

	"github.com/gin-gonic/gin"
)

// RespondError maps domain sentinel errors to HTTP status codes and the
// standard error envelope. It is the single place handlers translate service
// errors, so no raw SQL/GORM/stack details can leak to callers.
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
		// Never leak raw error text or DB internals to the API caller.
		Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
	}
}

// BindJSON binds and validates a JSON request body, writing a standard error
// envelope and returning false when binding fails.
func BindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return false
	}
	return true
}
