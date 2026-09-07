package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
