package user_test

import (
	"context"
	"testing"

	"i_m_s/internal/models"
	"i_m_s/internal/tests/mocks"
	"i_m_s/internal/user"
	authutils "i_m_s/internal/utils/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func validCreateInput() user.CreateUserInput {
	return user.CreateUserInput{
		FirstName: "Grace",
		LastName:  "Hopper",
		Username:  "grace",
		Email:     "grace@example.com",
		Password:  "plaintext-password",
		Role:      models.Staff,
	}
}

func TestService_CreateUserHashesPassword(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "grace").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "grace@example.com").Return(nil, models.ErrNotFound)

	var captured *models.User
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { captured = args.Get(2).(*models.User) }).
		Return(nil)

	created, err := user.NewService(nil, repo).CreateUser(context.Background(), validCreateInput())

	require.NoError(t, err)
	require.NotNil(t, captured)
	assert.NotEqual(t, "plaintext-password", captured.PassHash)
	assert.True(t, authutils.CheckPassword(captured.PassHash, "plaintext-password"))
	assert.True(t, created.IsActive)
	assert.Equal(t, models.Staff, created.Role)
}

func TestService_CreateUserRejectsDuplicateUsername(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "grace").Return(&models.User{}, nil)

	_, err := user.NewService(nil, repo).CreateUser(context.Background(), validCreateInput())

	assert.ErrorIs(t, err, models.ErrConflict)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateUserRejectsDuplicateEmail(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("GetByUsername", mock.Anything, mock.Anything, "grace").Return(nil, models.ErrNotFound)
	repo.On("GetByEmail", mock.Anything, mock.Anything, "grace@example.com").Return(&models.User{}, nil)

	_, err := user.NewService(nil, repo).CreateUser(context.Background(), validCreateInput())

	assert.ErrorIs(t, err, models.ErrConflict)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateUserRejectsInvalidRole(t *testing.T) {
	repo := new(mocks.UserRepository)
	input := validCreateInput()
	input.Role = models.Roles("superadmin")

	_, err := user.NewService(nil, repo).CreateUser(context.Background(), input)

	assert.ErrorIs(t, err, models.ErrBadRequest)
	repo.AssertNotCalled(t, "GetByUsername", mock.Anything, mock.Anything, mock.Anything)
}

func TestService_UpdateRole(t *testing.T) {
	repo := new(mocks.UserRepository)
	account := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	repo.On("GetByID", mock.Anything, mock.Anything, account.ID).Return(account, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	updated, err := user.NewService(nil, repo).UpdateRole(context.Background(), account.ID, models.Manager)

	require.NoError(t, err)
	assert.Equal(t, models.Manager, updated.Role)
}

func TestService_UpdateRoleRejectsInvalidRole(t *testing.T) {
	repo := new(mocks.UserRepository)

	_, err := user.NewService(nil, repo).UpdateRole(context.Background(), uuid.New(), models.Roles("owner"))

	assert.ErrorIs(t, err, models.ErrBadRequest)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
}

func TestService_UpdateRoleNotFound(t *testing.T) {
	repo := new(mocks.UserRepository)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, mock.Anything, id).Return(nil, models.ErrNotFound)

	_, err := user.NewService(nil, repo).UpdateRole(context.Background(), id, models.Manager)

	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestService_List(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("List", mock.Anything, mock.Anything, 1, 20).Return([]models.User{{}, {}}, int64(2), nil)

	users, total, err := user.NewService(nil, repo).List(context.Background(), 1, 20)

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, int64(2), total)
}
