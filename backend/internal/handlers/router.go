package handlers

import (
	"net/http"

	"i_m_s/internal/middleware"
	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter creates the Gin engine, registers global middleware,
// and configures the base route groups.
func SetupRouter(db *gorm.DB, jwtManager *authutils.JWTManager, isDev bool) *gin.Engine {
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.PanicRecovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())

	// Health check — lives outside /api/v1 so load balancers can hit it directly
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 route group
	v1 := router.Group("/api/v1")
	v1.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{"status": "ok", "version": "v1"})
	})

	// Domain routes (auth, users, ...)
	RegisterRoutes(v1, db, jwtManager)

	return router
}
