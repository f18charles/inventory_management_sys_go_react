package services_test

import (
	"context"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/services"
	authutils "i_m_s/internal/utils/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func validCreateInput() services.CreateUserInput {
	return services.CreateUserInput{
		FirstName: "Grace",
		LastName:  "Hopper",
		Username:  "grace",
		Email:     "grace@example.com",
		Password:  "plaintext-password",
		Role:      models.Staff,
	}
}

func TestUserService_CreateUserHashesPassword(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "grace").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "grace@example.com").Return(nil, models.ErrNotFound)

	var captured *models.User
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { captured = args.Get(2).(*models.User) }).
		Return(nil)

	user, err := services.NewUserService(nil, repo).CreateUser(context.Background(), validCreateInput())

	require.NoError(t, err)
	require.NotNil(t, captured)
	assert.NotEqual(t, "plaintext-password", captured.PassHash)
	assert.True(t, authutils.CheckPassword(captured.PassHash, "plaintext-password"))
	assert.True(t, user.IsActive)
	assert.Equal(t, models.Staff, user.Role)
}

func TestUserService_CreateUserRejectsDuplicateUsername(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "grace").Return(&models.User{}, nil)

	_, err := services.NewUserService(nil, repo).CreateUser(context.Background(), validCreateInput())

	assert.ErrorIs(t, err, models.ErrConflict)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_CreateUserRejectsDuplicateEmail(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "grace").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "grace@example.com").Return(&models.User{}, nil)

	_, err := services.NewUserService(nil, repo).CreateUser(context.Background(), validCreateInput())

	assert.ErrorIs(t, err, models.ErrConflict)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_CreateUserRejectsInvalidRole(t *testing.T) {
	repo := new(mockUserRepository)
	input := validCreateInput()
	input.Role = models.Roles("superadmin")

	_, err := services.NewUserService(nil, repo).CreateUser(context.Background(), input)

	assert.ErrorIs(t, err, models.ErrBadRequest)
	repo.AssertNotCalled(t, "GetByUsername", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_UpdateUserRole(t *testing.T) {
	repo := new(mockUserRepository)
	user := activeUser(t, "irrelevant")
	repo.On("GetByID", mock.Anything, mock.Anything, user.ID).Return(user, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	updated, err := services.NewUserService(nil, repo).UpdateUserRole(context.Background(), user.ID, models.Manager)

	require.NoError(t, err)
	assert.Equal(t, models.Manager, updated.Role)
}

func TestUserService_UpdateUserRoleRejectsInvalidRole(t *testing.T) {
	repo := new(mockUserRepository)

	_, err := services.NewUserService(nil, repo).UpdateUserRole(context.Background(), uuid.New(), models.Roles("owner"))

	assert.ErrorIs(t, err, models.ErrBadRequest)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_UpdateUserRoleNotFound(t *testing.T) {
	repo := new(mockUserRepository)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, mock.Anything, id).Return(nil, models.ErrNotFound)

	_, err := services.NewUserService(nil, repo).UpdateUserRole(context.Background(), id, models.Manager)

	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestUserService_ListUsers(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("List", mock.Anything, mock.Anything, 1, 20).Return([]models.User{{}, {}}, int64(2), nil)

	users, total, err := services.NewUserService(nil, repo).ListUsers(context.Background(), 1, 20)

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, int64(2), total)
}
