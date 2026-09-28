package mocks

import (
	"context"

	"i_m_s/internal/models"
	"i_m_s/internal/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// UserRepository is a testify-based test double for the user domain's
// Repository interface, shared by the auth and user domain tests.
type UserRepository struct {
	mock.Mock
}

var _ user.Repository = (*UserRepository)(nil)

func (m *UserRepository) Create(ctx context.Context, db *gorm.DB, u *models.User) error {
	args := m.Called(ctx, db, u)
	return args.Error(0)
}

func (m *UserRepository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) GetByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error) {
	args := m.Called(ctx, db, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) GetByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error) {
	args := m.Called(ctx, db, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) Update(ctx context.Context, db *gorm.DB, u *models.User) error {
	args := m.Called(ctx, db, u)
	return args.Error(0)
}

func (m *UserRepository) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *UserRepository) List(ctx context.Context, db *gorm.DB, page, pageSize int) ([]models.User, int64, error) {
	args := m.Called(ctx, db, page, pageSize)
	users, _ := args.Get(0).([]models.User)
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return users, total, args.Error(2)
}
