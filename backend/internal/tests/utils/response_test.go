package utils_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewPaginationMeta(t *testing.T) {
	meta := response.NewPaginationMeta(1, 20, 45)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 20, meta.PageSize)
	assert.Equal(t, int64(45), meta.TotalItems)
	assert.Equal(t, 3, meta.TotalPages) // 45 / 20 ceiling = 3

	metaDefault := response.NewPaginationMeta(1, 0, 10)
	assert.Equal(t, 20, metaDefault.PageSize)
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	payload := map[string]string{"message": "hello world"}
	response.Success(c, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data map[string]string `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "hello world", res.Data["message"])
}

func TestPaginatedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []string{"item1", "item2"}
	meta := response.NewPaginationMeta(1, 10, 2)
	response.Paginated(c, http.StatusOK, items, meta)

	assert.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data []string               `json:"data"`
		Meta response.PaginationMeta `json:"meta"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, 1, res.Meta.Page)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
}

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "field X is required")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res struct {
		Error response.ErrorBody `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "INVALID_INPUT", res.Error.Code)
	assert.Equal(t, "field X is required", res.Error.Message)
}

func TestRespondErrorMapping(t *testing.T) {
	cases := []struct {
		err        error
		wantStatus int
		wantCode   string
	}{
		{models.ErrNotFound, http.StatusNotFound, "NOT_FOUND"},
		{models.ErrInsufficientInventory, http.StatusConflict, "INSUFFICIENT_INVENTORY"},
		{models.ErrUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{models.ErrInvalidCredentials, http.StatusUnauthorized, "INVALID_CREDENTIALS"},
		{models.ErrAccountInactive, http.StatusForbidden, "ACCOUNT_INACTIVE"},
		{models.ErrForbidden, http.StatusForbidden, "FORBIDDEN"},
		{models.ErrBadRequest, http.StatusBadRequest, "BAD_REQUEST"},
		{models.ErrConflict, http.StatusConflict, "CONFLICT"},
		{models.ErrValidationError, http.StatusUnprocessableEntity, "VALIDATION_ERROR"},
		{fmt.Errorf("some raw db error"), http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		response.RespondError(c, tc.err)

		assert.Equal(t, tc.wantStatus, w.Code)

		var res struct {
			Error response.ErrorBody `json:"error"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, tc.wantCode, res.Error.Code)
	}
}

func TestBindJSON(t *testing.T) {
	r := gin.New()
	r.POST("/test-bind", func(c *gin.Context) {
		var req struct {
			Name string `json:"name" binding:"required"`
		}
		if response.BindJSON(c, &req) {
			response.Success(c, http.StatusOK, req)
		}
	})

	wValid := httptest.NewRecorder()
	reqValid, _ := http.NewRequest(http.MethodPost, "/test-bind", strings.NewReader(`{"name":"Product A"}`))
	reqValid.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wValid, reqValid)
	assert.Equal(t, http.StatusOK, wValid.Code)

	wInvalid := httptest.NewRecorder()
	reqInvalid, _ := http.NewRequest(http.MethodPost, "/test-bind", strings.NewReader(`{}`))
	reqInvalid.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wInvalid, reqInvalid)
	assert.Equal(t, http.StatusBadRequest, wInvalid.Code)
}
