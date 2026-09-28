package user_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/tests/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ListAsAdmin(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("List", mock.Anything, mock.Anything, 1, 20).
		Return([]models.User{{BaseModel: models.BaseModel{ID: uuid.New()}, Username: "a"}}, int64(1), nil)
	engine, jwtManager := newUserRouter(repo)

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodGet, "/api/v1/users", "", token)

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data []models.User `json:"data"`
		Meta struct {
			TotalItems int64 `json:"total_items"`
		} `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Len(t, res.Data, 1)
	assert.Equal(t, int64(1), res.Meta.TotalItems)
}

func TestHandler_ListForbiddenForStaff(t *testing.T) {
	engine, jwtManager := newUserRouter(new(mocks.UserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Staff)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodGet, "/api/v1/users", "", token)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_ListRequiresAuth(t *testing.T) {
	engine, _ := newUserRouter(new(mocks.UserRepository))

	w := doRequest(t, engine, http.MethodGet, "/api/v1/users", "", "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_CreateAsAdmin(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "newbie").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "newbie@example.com").Return(nil, models.ErrNotFound)
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	engine, jwtManager := newUserRouter(repo)

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	body := `{"first_name":"New","last_name":"Bie","username":"newbie","email":"newbie@example.com","password":"password123","role":"staff"}`
	w := doRequest(t, engine, http.MethodPost, "/api/v1/users", body, token)

	require.Equal(t, http.StatusCreated, w.Code)

	var res struct {
		Data models.User `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, "newbie", res.Data.Username)
	assert.NotContains(t, w.Body.String(), "password123")
}

func TestHandler_CreateForbiddenForStaff(t *testing.T) {
	engine, jwtManager := newUserRouter(new(mocks.UserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Staff)
	require.NoError(t, err)

	body := `{"first_name":"New","last_name":"Bie","username":"newbie","email":"newbie@example.com","password":"password123","role":"staff"}`
	w := doRequest(t, engine, http.MethodPost, "/api/v1/users", body, token)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_CreateInvalidBody(t *testing.T) {
	engine, jwtManager := newUserRouter(new(mocks.UserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodPost, "/api/v1/users",
		`{"first_name":"New","last_name":"Bie","username":"newbie","email":"not-an-email","password":"password123","role":"staff"}`, token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_UpdateRoleAsAdmin(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}, Role: models.Staff}
	repo.On("GetByID", mock.Anything, mock.Anything, account.ID).Return(account, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	engine, jwtManager := newUserRouter(repo)

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodPatch, "/api/v1/users/"+account.ID.String()+"/role", `{"role":"manager"}`, token)

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data models.User `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, models.Manager, res.Data.Role)
}

func TestHandler_UpdateRoleInvalidID(t *testing.T) {
	engine, jwtManager := newUserRouter(new(mocks.UserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodPatch, "/api/v1/users/not-a-uuid/role", `{"role":"manager"}`, token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
