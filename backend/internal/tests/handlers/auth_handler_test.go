package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"i_m_s/internal/models"
	authutils "i_m_s/internal/utils/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func activeHandlerUser(t *testing.T, password string) *models.User {
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

func TestAuthHandler_LoginSuccess(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeHandlerUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(user, nil)
	engine, _ := newTestRouter(repo)

	w := doRequest(t, engine, http.MethodPost, "/api/v1/auth/login",
		`{"username_or_email":"ada","password":"correct-password"}`, "")

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.NotEmpty(t, res.Data.Token)
	assert.Equal(t, "ada", res.Data.User.Username)

	body := w.Body.String()
	assert.NotContains(t, body, "pass_hash")
	assert.NotContains(t, body, user.PassHash)
}

func TestAuthHandler_LoginInvalidBody(t *testing.T) {
	engine, _ := newTestRouter(new(mockUserRepository))

	w := doRequest(t, engine, http.MethodPost, "/api/v1/auth/login", `{"username_or_email":"ada"}`, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_LoginWrongPassword(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeHandlerUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(user, nil)
	engine, _ := newTestRouter(repo)

	w := doRequest(t, engine, http.MethodPost, "/api/v1/auth/login",
		`{"username_or_email":"ada","password":"wrong"}`, "")

	require.Equal(t, http.StatusUnauthorized, w.Code)
	assertBodyCode(t, w, "INVALID_CREDENTIALS")
	assert.False(t, strings.Contains(w.Body.String(), "correct-password"))
}

func TestAuthHandler_LoginInactiveAccount(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeHandlerUser(t, "correct-password")
	user.IsActive = false
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(user, nil)
	engine, _ := newTestRouter(repo)

	w := doRequest(t, engine, http.MethodPost, "/api/v1/auth/login",
		`{"username_or_email":"ada","password":"correct-password"}`, "")

	require.Equal(t, http.StatusForbidden, w.Code)
	assertBodyCode(t, w, "ACCOUNT_INACTIVE")
}

func TestAuthHandler_MeRequiresToken(t *testing.T) {
	engine, _ := newTestRouter(new(mockUserRepository))

	w := doRequest(t, engine, http.MethodGet, "/api/v1/auth/me", "", "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_MeReturnsProfile(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeHandlerUser(t, "correct-password")
	repo.On("GetByID", mock.Anything, mock.Anything, user.ID).Return(user, nil)
	engine, jwtManager := newTestRouter(repo)

	token, err := jwtManager.Generate(user.ID, models.Staff)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodGet, "/api/v1/auth/me", "", token)

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data models.User `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, user.ID, res.Data.ID)
	assert.NotContains(t, w.Body.String(), "pass_hash")
}
