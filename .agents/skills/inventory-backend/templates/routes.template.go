// internal/user/routes.go
//
// Routes are OWNED by the domain and live in its package. They mount the
// handler on the /api/v1 group and attach JWTAuthMiddleware + RequireRole.

package user

import (
	"github.com/gin-gonic/gin"

	"i_m_s/internal/middleware"
	"i_m_s/internal/models"
	authutils "i_m_s/internal/utils/auth"
)

func RegisterRoutes(v1 *gin.RouterGroup, handler *Handler, jwtManager *authutils.JWTManager) {
	users := v1.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(jwtManager))
	{
		users.GET("", middleware.RequireRole(models.Admin, models.Manager), handler.List)
		users.POST("", middleware.RequireRole(models.Admin), handler.Create)
		users.PATCH("/:id/role", middleware.RequireRole(models.Admin), handler.UpdateRole)
	}
}

// ---------------------------------------------------------------------------
// internal/router/router.go
//
// The router package owns Gin engine setup and dependency injection. It
// constructs each domain's repository/service/handler and mounts its routes;
// domains never import another domain's handler or the router itself.
//
// func SetupRouter(db *gorm.DB, jwtManager *authutils.JWTManager, isDev bool) *gin.Engine {
//     if !isDev { gin.SetMode(gin.ReleaseMode) }
//     engine := gin.New()
//     engine.Use(middleware.PanicRecovery(), middleware.CORS(), middleware.RequestLogger())
//     engine.GET("/health", func(c *gin.Context) { response.Success(c, http.StatusOK, gin.H{"status": "ok"}) })
//     v1 := engine.Group("/api/v1")
//     v1.GET("/health", func(c *gin.Context) { response.Success(c, http.StatusOK, gin.H{"status": "ok", "version": "v1"}) })
//     RegisterDomains(v1, db, jwtManager)
//     return engine
// }
//
// func RegisterDomains(v1 *gin.RouterGroup, db *gorm.DB, jwtManager *authutils.JWTManager) {
//     userRepo := user.NewRepository()
//     user.RegisterRoutes(v1, user.NewHandler(user.NewService(db, userRepo)), jwtManager)
//     // auth depends on the user domain's repository (auth -> user).
//     auth.RegisterRoutes(v1, auth.NewHandler(auth.NewService(db, userRepo, jwtManager)), jwtManager)
// }
