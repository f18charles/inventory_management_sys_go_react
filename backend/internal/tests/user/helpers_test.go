package user_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"i_m_s/internal/tests/mocks"
	"i_m_s/internal/user"
	authutils "i_m_s/internal/utils/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newUserRouter(repo *mocks.UserRepository) (*gin.Engine, *authutils.JWTManager) {
	jwtManager := authutils.NewJWTManager("test-secret", time.Hour)
	engine := gin.New()
	v1 := engine.Group("/api/v1")
	user.RegisterRoutes(v1, user.NewHandler(user.NewService(nil, repo)), jwtManager)
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
