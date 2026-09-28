package middleware

import (
	"net/http"

	"i_m_s/internal/models"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// RequireRole allows the request through only when the authenticated user's
// role is one of the supplied roles. It must run after JWTAuthMiddleware.
func RequireRole(allowed ...models.Roles) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(ContextUserRole)
		if !exists {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "access denied")
			c.Abort()
			return
		}

		for _, candidate := range allowed {
			if models.Roles(role) == candidate {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "FORBIDDEN", "insufficient permissions")
		c.Abort()
	}
}
