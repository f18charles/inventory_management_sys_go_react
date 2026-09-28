package user

import (
	"i_m_s/internal/middleware"
	"i_m_s/internal/models"
	authutils "i_m_s/internal/utils/auth"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires the user management endpoints behind authentication and
// role-based authorization.
func RegisterRoutes(v1 *gin.RouterGroup, handler *Handler, jwtManager *authutils.JWTManager) {
	users := v1.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(jwtManager))
	{
		users.GET("", middleware.RequireRole(models.Admin, models.Manager), handler.List)
		users.POST("", middleware.RequireRole(models.Admin), handler.Create)
		users.PATCH("/:id/role", middleware.RequireRole(models.Admin), handler.UpdateRole)
	}
}
