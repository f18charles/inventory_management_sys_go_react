package router

import (
	"net/http"

	"i_m_s/internal/auth"
	"i_m_s/internal/middleware"
	"i_m_s/internal/user"
	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter creates the Gin engine, registers global middleware, and wires
// every domain's routes onto the /api/v1 group.
func SetupRouter(db *gorm.DB, jwtManager *authutils.JWTManager, isDev bool) *gin.Engine {
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// Global middleware
	engine.Use(middleware.PanicRecovery())
	engine.Use(middleware.CORS())
	engine.Use(middleware.RequestLogger())

	// Health checks — outside /api/v1 so load balancers can hit them directly
	engine.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := engine.Group("/api/v1")
	v1.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{"status": "ok", "version": "v1"})
	})

	RegisterDomains(v1, db, jwtManager)

	return engine
}

// RegisterDomains performs dependency injection for each domain and mounts its
// routes. Repositories are shared where a domain depends on another
// (auth -> user).
func RegisterDomains(v1 *gin.RouterGroup, db *gorm.DB, jwtManager *authutils.JWTManager) {
	userRepo := user.NewRepository()

	auth.RegisterRoutes(v1, auth.NewHandler(auth.NewService(db, userRepo, jwtManager)), jwtManager)
	user.RegisterRoutes(v1, user.NewHandler(user.NewService(db, userRepo)), jwtManager)
}
