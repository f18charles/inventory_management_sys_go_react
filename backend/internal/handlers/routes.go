package handlers

import (
	"i_m_s/internal/middleware"
	"i_m_s/internal/models"
	"i_m_s/internal/repositories"
	"i_m_s/internal/services"
	authutils "i_m_s/internal/utils/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAuthRoutes wires the authentication endpoints.
func RegisterAuthRoutes(v1 *gin.RouterGroup, handler *AuthHandler, jwtManager *authutils.JWTManager) {
	auth := v1.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.GET("/me", middleware.JWTAuthMiddleware(jwtManager), handler.Me)
	}
}

// RegisterUserRoutes wires the user management endpoints behind auth and RBAC.
func RegisterUserRoutes(v1 *gin.RouterGroup, handler *UserHandler, jwtManager *authutils.JWTManager) {
	users := v1.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(jwtManager))
	{
		users.GET("", middleware.RequireRole(models.Admin, models.Manager), handler.List)
		users.POST("", middleware.RequireRole(models.Admin), handler.Create)
		users.PATCH("/:id/role", middleware.RequireRole(models.Admin), handler.UpdateRole)
	}
}

// RegisterRoutes composes the repositories and services for the auth/user
// domains and registers their routes on the supplied v1 group.
func RegisterRoutes(v1 *gin.RouterGroup, db *gorm.DB, jwtManager *authutils.JWTManager) {
	userRepo := repositories.NewUserRepository()

	RegisterAuthRoutes(v1, NewAuthHandler(services.NewAuthService(db, userRepo, jwtManager)), jwtManager)
	RegisterUserRoutes(v1, NewUserHandler(services.NewUserService(db, userRepo)), jwtManager)
}
