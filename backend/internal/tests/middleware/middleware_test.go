package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"i_m_s/internal/middleware"
	"i_m_s/internal/utils/logger"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
	logger.Init("development")
}

func TestCORSMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS())
	r.GET("/test-cors", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/test-cors", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestRequestLoggerMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.GET("/test-logger", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-logger", nil)
	req.Header.Set("X-Request-ID", "custom-req-id-123")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "custom-req-id-123", w.Header().Get("X-Request-ID"))
}

func TestPanicRecoveryMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(middleware.PanicRecovery())
	r.GET("/test-panic", func(c *gin.Context) {
		panic("something went critically wrong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-panic", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res struct {
		Error response.ErrorBody `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", res.Error.Code)
	assert.Equal(t, "An unexpected error occurred", res.Error.Message)
}
