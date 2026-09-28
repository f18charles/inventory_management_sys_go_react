package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"i_m_s/internal/middleware"
	"i_m_s/internal/models"
	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func protectedRouter(jwtManager *authutils.JWTManager) *gin.Engine {
	r := gin.New()
	r.Use(middleware.JWTAuthMiddleware(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"user_id": c.GetString(middleware.ContextUserID),
			"role":    c.GetString(middleware.ContextUserRole),
		})
	})
	return r
}

func TestJWTAuthMiddleware_RejectsMissingHeader(t *testing.T) {
	r := protectedRouter(authutils.NewJWTManager("secret", time.Hour))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_RejectsMalformedHeader(t *testing.T) {
	r := protectedRouter(authutils.NewJWTManager("secret", time.Hour))

	for _, header := range []string{"Token abc", "Bearer", "Bearer ", "abc"} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", header)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "header=%q", header)
	}
}

func TestJWTAuthMiddleware_RejectsInvalidToken(t *testing.T) {
	r := protectedRouter(authutils.NewJWTManager("secret", time.Hour))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_RejectsExpiredToken(t *testing.T) {
	mgr := authutils.NewJWTManager("secret", -time.Minute)
	token, err := mgr.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	protectedRouter(mgr).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_AttachesIdentity(t *testing.T) {
	mgr := authutils.NewJWTManager("secret", time.Hour)
	userID := uuid.New()
	token, err := mgr.Generate(userID, models.Manager)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	protectedRouter(mgr).ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data struct {
			UserID string `json:"user_id"`
			Role   string `json:"role"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, userID.String(), res.Data.UserID)
	assert.Equal(t, "manager", res.Data.Role)
}

func roleRouter(role string, allowed ...models.Roles) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if role != "" {
			c.Set(middleware.ContextUserRole, role)
		}
		c.Next()
	})
	r.GET("/restricted", middleware.RequireRole(allowed...), func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestRequireRole(t *testing.T) {
	cases := []struct {
		name string
		role string
		want int
	}{
		{"admin is allowed", "admin", http.StatusOK},
		{"manager is allowed", "manager", http.StatusOK},
		{"staff is forbidden", "staff", http.StatusForbidden},
		{"anonymous is unauthorized", "", http.StatusUnauthorized},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/restricted", nil)
			roleRouter(tc.role, models.Admin, models.Manager).ServeHTTP(w, req)

			assert.Equal(t, tc.want, w.Code)
		})
	}
}

func TestRequireRole_RejectsUnknownRoleValue(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/restricted", nil)
	roleRouter("owner", models.Admin).ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
