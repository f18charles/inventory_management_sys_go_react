package auth_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_LoginSuccess(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(account, nil)
	engine, _ := newAuthRouter(repo)

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
	assert.NotContains(t, body, account.PassHash)
}

func TestHandler_LoginInvalidBody(t *testing.T) {
	engine, _ := newAuthRouter(new(mocks.UserRepository))

	w := doRequest(t, engine, http.MethodPost, "/api/v1/auth/login", `{"username_or_email":"ada"}`, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_LoginWrongPassword(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(account, nil)
	engine, _ := newAuthRouter(repo)

	w := doRequest(t, engine, http.MethodPost, "/api/v1/auth/login",
		`{"username_or_email":"ada","password":"wrong"}`, "")

	require.Equal(t, http.StatusUnauthorized, w.Code)
	assertBodyCode(t, w, "INVALID_CREDENTIALS")
	assert.False(t, strings.Contains(w.Body.String(), "correct-password"))
}

func TestHandler_LoginInactiveAccount(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	account.IsActive = false
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(account, nil)
	engine, _ := newAuthRouter(repo)

	w := doRequest(t, engine, http.MethodPost, "/api/v1/auth/login",
		`{"username_or_email":"ada","password":"correct-password"}`, "")

	require.Equal(t, http.StatusForbidden, w.Code)
	assertBodyCode(t, w, "ACCOUNT_INACTIVE")
}

func TestHandler_MeRequiresToken(t *testing.T) {
	engine, _ := newAuthRouter(new(mocks.UserRepository))

	w := doRequest(t, engine, http.MethodGet, "/api/v1/auth/me", "", "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_MeReturnsProfile(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	repo.On("GetByID", mock.Anything, mock.Anything, account.ID).Return(account, nil)
	engine, jwtManager := newAuthRouter(repo)

	token, err := jwtManager.Generate(account.ID, models.Staff)
	require.NoError(t, err)

	w := doRequest(t, engine, http.MethodGet, "/api/v1/auth/me", "", token)

	require.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Data models.User `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, account.ID, res.Data.ID)
	assert.NotContains(t, w.Body.String(), "pass_hash")
}
