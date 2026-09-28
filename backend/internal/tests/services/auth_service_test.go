package services_test

import (
	"context"
	"testing"
	"time"

	"i_m_s/internal/models"
	"i_m_s/internal/services"
	authutils "i_m_s/internal/utils/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func activeUser(t *testing.T, password string) *models.User {
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

func newAuthService(repo *mockUserRepository) *services.AuthService {
	return services.NewAuthService(nil, repo, authutils.NewJWTManager("test-secret", time.Hour))
}

func TestAuthService_LoginByUsername(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(user, nil)

	got, token, err := newAuthService(repo).Login(context.Background(), "ada", "correct-password")

	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.NotEmpty(t, token)
}

func TestAuthService_LoginByEmailFallback(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada@example.com").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "ada@example.com").Return(user, nil)

	got, token, err := newAuthService(repo).Login(context.Background(), "ada@example.com", "correct-password")

	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.NotEmpty(t, token)
}

func TestAuthService_LoginRejectsWrongPassword(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeUser(t, "correct-password")
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(user, nil)

	_, _, err := newAuthService(repo).Login(context.Background(), "ada", "wrong-password")

	assert.ErrorIs(t, err, models.ErrInvalidCredentials)
}

func TestAuthService_LoginRejectsUnknownUser(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ghost").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "ghost").Return(nil, models.ErrNotFound)

	_, _, err := newAuthService(repo).Login(context.Background(), "ghost", "whatever")

	assert.ErrorIs(t, err, models.ErrInvalidCredentials)
}

func TestAuthService_LoginRejectsInactiveAccount(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeUser(t, "correct-password")
	user.IsActive = false
	repo.On("GetByUsername", mock.Anything, mock.Anything, "ada").Return(user, nil)

	_, _, err := newAuthService(repo).Login(context.Background(), "ada", "correct-password")

	assert.ErrorIs(t, err, models.ErrAccountInactive)
}

func TestAuthService_GetProfile(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeUser(t, "correct-password")
	repo.On("GetByID", mock.Anything, mock.Anything, user.ID).Return(user, nil)

	got, err := newAuthService(repo).GetProfile(context.Background(), user.ID)

	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}
