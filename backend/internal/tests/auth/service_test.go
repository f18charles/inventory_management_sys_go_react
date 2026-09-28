package auth_test

import (
	"context"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_LoginByUsername(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(account, nil)

	got, token, err := newAuthService(repo).Login(context.Background(), "ada", "correct-password")

	require.NoError(t, err)
	assert.Equal(t, account.ID, got.ID)
	assert.NotEmpty(t, token)
}

func TestService_LoginByEmailFallback(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada@example.com").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "ada@example.com").Return(account, nil)

	got, token, err := newAuthService(repo).Login(context.Background(), "ada@example.com", "correct-password")

	require.NoError(t, err)
	assert.Equal(t, account.ID, got.ID)
	assert.NotEmpty(t, token)
}

func TestService_LoginRejectsWrongPassword(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(account, nil)

	_, _, err := newAuthService(repo).Login(context.Background(), "ada", "wrong-password")

	assert.ErrorIs(t, err, models.ErrInvalidCredentials)
}

func TestService_LoginRejectsUnknownUser(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ghost").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "ghost").Return(nil, models.ErrNotFound)

	_, _, err := newAuthService(repo).Login(context.Background(), "ghost", "whatever")

	assert.ErrorIs(t, err, models.ErrInvalidCredentials)
}

func TestService_LoginRejectsInactiveAccount(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	account.IsActive = false
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(account, nil)

	_, _, err := newAuthService(repo).Login(context.Background(), "ada", "correct-password")

	assert.ErrorIs(t, err, models.ErrAccountInactive)
}

func TestService_GetProfile(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := activeUser(t, "correct-password")
	repo.On("GetByID", mock.Anything, mock.Anything, account.ID).Return(account, nil)

	got, err := newAuthService(repo).GetProfile(context.Background(), account.ID)

	require.NoError(t, err)
	assert.Equal(t, account.ID, got.ID)
}
