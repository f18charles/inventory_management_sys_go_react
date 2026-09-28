package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"i_m_s/internal/auth"
	"i_m_s/internal/models"
	"i_m_s/internal/tests/mocks"
	authutils "i_m_s/internal/utils/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func activeUser(t *testing.T, password string) *models.User {
	t.Helper()
	hash, err := authutils.HashPassword(password)
	require.NoError(t, err)
	return &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		FirstName: "Ada",
		LastName:  "Lovelace",
		Username:  "ada",
		Email:     "ada@example.com",
		PassHash:  hash,
		Role:      models.Staff,
		IsActive:  true,
	}
}

func newAuthService(repo *mocks.UserRepository) *auth.Service {
	return auth.NewService(nil, repo, authutils.NewJWTManager("test-secret", time.Hour))
}

func newAuthRouter(repo *mocks.UserRepository) (*gin.Engine, *authutils.JWTManager) {
	jwtManager := authutils.NewJWTManager("test-secret", time.Hour)
	engine := gin.New()
	v1 := engine.Group("/api/v1")
	auth.RegisterRoutes(v1, auth.NewHandler(auth.NewService(nil, repo, jwtManager)), jwtManager)
	return engine, jwtManager
}

func doRequest(t *testing.T, engine *gin.Engine, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, path, reader)
	require.NoError(t, err)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func assertBodyCode(t *testing.T, w *httptest.ResponseRecorder, code string) {
	t.Helper()

	var res struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, code, res.Error.Code)
}
