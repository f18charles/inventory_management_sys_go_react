package repositories

import (
	"context"
	"errors"
	"fmt"

	"i_m_s/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository is the persistence boundary for user records. It exposes pure
// CRUD/query operations and contains no business or authorization decisions.
type UserRepository interface {
	Create(ctx context.Context, db *gorm.DB, user *models.User) error
	GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error)
	GetByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error)
	Update(ctx context.Context, db *gorm.DB, user *models.User) error
	Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error
	List(ctx context.Context, db *gorm.DB, page, pageSize int) ([]models.User, int64, error)
}

type userRepository struct{}

// NewUserRepository returns a GORM-backed UserRepository. It accepts a
// transaction handle per call so the same repository works standalone or
// inside a service-owned transaction.
func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user %s: %w", id, err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := db.WithContext(ctx).First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.WithContext(ctx).First(&user, "username = ?", username).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, db *gorm.DB, user *models.User) error {
	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("update user %s: %w", user.ID, err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	res := db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete user %s: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, db *gorm.DB, page, pageSize int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var total int64
	if err := db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	var users []models.User
	offset := (page - 1) * pageSize
	if err := db.WithContext(ctx).
		Model(&models.User{}).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return users, total, nil
}
