package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"i_m_s/internal/models"
	"i_m_s/internal/utils/response"
)

// RespondError maps domain sentinel errors to HTTP response envelopes.
func RespondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "resource not found")
	case errors.Is(err, models.ErrInsufficientInventory):
		response.Error(c, http.StatusConflict, "INSUFFICIENT_INVENTORY", err.Error())
	case errors.Is(err, models.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	case errors.Is(err, models.ErrForbidden):
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "forbidden")
	case errors.Is(err, models.ErrBadRequest):
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	case errors.Is(err, models.ErrConflict):
		response.Error(c, http.StatusConflict, "CONFLICT", err.Error())
	case errors.Is(err, models.ErrValidationError):
		response.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
	default:
		// Never leak raw error text or DB internals to the API caller
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
	}
}

// BindAndValidate binds JSON request bodies and handles validation failure using standard error envelopes.
func BindAndValidate(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return false
	}
	return true
}
