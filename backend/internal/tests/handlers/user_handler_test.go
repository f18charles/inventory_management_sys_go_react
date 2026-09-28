package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/utils/response"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_ListAsAdmin(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("List", mock.Anything, mock.Anything, 1, 20).
		Return([]models.User{{BaseModel: models.BaseModel{ID: uuid.New()}, Username: "a"}}, int64(1), nil)
	engine, jwtManager := newTestRouter(repo)

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodGet, "/api/v1/users", "", token)

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data []models.User           `json:"data"`
		Meta response.PaginationMeta `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Len(t, res.Data, 1)
	assert.Equal(t, int64(1), res.Meta.TotalItems)
}

func TestUserHandler_ListForbiddenForStaff(t *testing.T) {
	engine, jwtManager := newTestRouter(new(mockUserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Staff)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodGet, "/api/v1/users", "", token)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserHandler_ListRequiresAuth(t *testing.T) {
	engine, _ := newTestRouter(new(mockUserRepository))

	w := doRequest(t, engine, http.MethodGet, "/api/v1/users", "", "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_CreateAsAdmin(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "newbie").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "newbie@example.com").Return(nil, models.ErrNotFound)
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	engine, jwtManager := newTestRouter(repo)

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

func TestUserHandler_CreateForbiddenForStaff(t *testing.T) {
	engine, jwtManager := newTestRouter(new(mockUserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Staff)
	require.NoError(t, err)

	body := `{"first_name":"New","last_name":"Bie","username":"newbie","email":"newbie@example.com","password":"password123","role":"staff"}`
	w := doRequest(t, engine, http.MethodPost, "/api/v1/users", body, token)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserHandler_CreateInvalidBody(t *testing.T) {
	engine, jwtManager := newTestRouter(new(mockUserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodPost, "/api/v1/users",
		`{"first_name":"New","last_name":"Bie","username":"newbie","email":"not-an-email","password":"password123","role":"staff"}`, token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UpdateRoleAsAdmin(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeHandlerUser(t, "correct-password")
	repo.On("GetByID", mock.Anything, mock.Anything, user.ID).Return(user, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	engine, jwtManager := newTestRouter(repo)

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodPatch, "/api/v1/users/"+user.ID.String()+"/role", `{"role":"manager"}`, token)

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data models.User `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, models.Manager, res.Data.Role)
}

func TestUserHandler_UpdateRoleInvalidID(t *testing.T) {
	engine, jwtManager := newTestRouter(new(mockUserRepository))

	token, err := jwtManager.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodPatch, "/api/v1/users/not-a-uuid/role", `{"role":"manager"}`, token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
