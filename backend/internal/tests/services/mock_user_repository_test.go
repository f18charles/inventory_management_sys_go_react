package services_test

import (
	"context"

	"i_m_s/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// mockUserRepository is an in-memory test double for repositories.UserRepository.
// It records calls and returns programmed results so service business rules can
// be tested without a database.
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	args := m.Called(ctx, db, user)
	return args.Error(0)
}

func (m *mockUserRepository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error) {
	args := m.Called(ctx, db, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error) {
	args := m.Called(ctx, db, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, db *gorm.DB, user *models.User) error {
	args := m.Called(ctx, db, user)
	return args.Error(0)
}

func (m *mockUserRepository) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *mockUserRepository) List(ctx context.Context, db *gorm.DB, page, pageSize int) ([]models.User, int64, error) {
	args := m.Called(ctx, db, page, pageSize)
	users, _ := args.Get(0).([]models.User)
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return users, total, args.Error(2)
}
