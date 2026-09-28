package auth

import (
	"i_m_s/internal/middleware"
	authutils "i_m_s/internal/utils/auth"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires the authentication endpoints onto the v1 group.
func RegisterRoutes(v1 *gin.RouterGroup, handler *Handler, jwtManager *authutils.JWTManager) {
	auth := v1.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.GET("/me", middleware.JWTAuthMiddleware(jwtManager), handler.Me)
	}
}
