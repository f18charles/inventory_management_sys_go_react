package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"i_m_s/internal/handlers"
	"i_m_s/internal/models"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthCheck(t *testing.T) {
	router := handlers.SetupRouter(nil, true)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data map[string]string `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "ok", res.Data["status"])
}

func TestApiV1HealthCheck(t *testing.T) {
	router := handlers.SetupRouter(nil, true)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data map[string]string `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "ok", res.Data["status"])
	assert.Equal(t, "v1", res.Data["version"])
}

func TestRespondError(t *testing.T) {
	testCases := []struct {
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{models.ErrNotFound, http.StatusNotFound, "NOT_FOUND"},
		{models.ErrInsufficientInventory, http.StatusConflict, "INSUFFICIENT_INVENTORY"},
		{models.ErrUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{models.ErrForbidden, http.StatusForbidden, "FORBIDDEN"},
		{models.ErrBadRequest, http.StatusBadRequest, "BAD_REQUEST"},
		{models.ErrConflict, http.StatusConflict, "CONFLICT"},
		{models.ErrValidationError, http.StatusUnprocessableEntity, "VALIDATION_ERROR"},
		{fmt.Errorf("some raw db error"), http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handlers.RespondError(c, tc.err)

		assert.Equal(t, tc.expectedCode, w.Code)

		var res struct {
			Error response.ErrorBody `json:"error"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedMsg, res.Error.Code)
	}
}

type sampleReq struct {
	Name string `json:"name" binding:"required"`
}

func TestBindAndValidate(t *testing.T) {
	r := gin.New()
	r.POST("/test-bind", func(c *gin.Context) {
		var req sampleReq
		if handlers.BindAndValidate(c, &req) {
			response.Success(c, http.StatusOK, req)
		}
	})

	// Case 1: Valid payload
	wValid := httptest.NewRecorder()
	reqValid, _ := http.NewRequest(http.MethodPost, "/test-bind", bytes.NewBufferString(`{"name":"Product A"}`))
	reqValid.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wValid, reqValid)

	assert.Equal(t, http.StatusOK, wValid.Code)

	// Case 2: Invalid payload
	wInvalid := httptest.NewRecorder()
	reqInvalid, _ := http.NewRequest(http.MethodPost, "/test-bind", bytes.NewBufferString(`{}`))
	reqInvalid.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wInvalid, reqInvalid)

	assert.Equal(t, http.StatusBadRequest, wInvalid.Code)
}
