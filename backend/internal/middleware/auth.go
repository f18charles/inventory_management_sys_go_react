package middleware

import (
	"net/http"
	"strings"

	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/logger"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// Gin context keys set by JWTAuthMiddleware and read by RequireRole/handlers.
const (
	ContextUserID   = "user_id"
	ContextUserRole = "user_role"
)

// JWTAuthMiddleware validates the Bearer token and attaches the authenticated
// user's id and role to the request context.
func JWTAuthMiddleware(jwtManager *authutils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtManager.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			logger.LogError(err, "jwt authentication failed", logger.Fields{
				"path": c.Request.URL.Path,
			})
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserRole, string(claims.Role))
		c.Next()
	}
}
