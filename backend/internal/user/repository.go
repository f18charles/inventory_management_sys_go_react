package user

import (
	"context"
	"errors"
	"fmt"

	"i_m_s/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository is the persistence boundary for user records. It exposes pure
// CRUD/query operations and contains no business or authorization decisions.
// It accepts a *gorm.DB per call so it works standalone or inside a caller's
// transaction.
type Repository interface {
	Create(ctx context.Context, db *gorm.DB, user *models.User) error
	GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error)
	GetByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error)
	Update(ctx context.Context, db *gorm.DB, user *models.User) error
	Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error
	List(ctx context.Context, db *gorm.DB, page, pageSize int) ([]models.User, int64, error)
}

type repository struct{}

// NewRepository returns a GORM-backed User user.Repository.
func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	err := db.WithContext(ctx).Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return models.ErrConflict
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.User, error) {
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

func (r *repository) GetByEmail(ctx context.Context, db *gorm.DB, email string) (*models.User, error) {
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

func (r *repository) GetByUsername(ctx context.Context, db *gorm.DB, username string) (*models.User, error) {
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

func (r *repository) Update(ctx context.Context, db *gorm.DB, user *models.User) error {
	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("update user %s: %w", user.ID, err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	res := db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete user %s: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *repository) List(ctx context.Context, db *gorm.DB, page, pageSize int) ([]models.User, int64, error) {
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
