package services

import (
	"context"
	"errors"
	"fmt"

	"i_m_s/internal/models"
	"i_m_s/internal/repositories"
	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// UserService owns user lifecycle business rules: uniqueness, role validity,
// and password hashing. It depends on the repository interface, never GORM
// directly for query logic.
type UserService struct {
	db       *gorm.DB
	userRepo repositories.UserRepository
}

func NewUserService(db *gorm.DB, userRepo repositories.UserRepository) *UserService {
	return &UserService{db: db, userRepo: userRepo}
}

// CreateUserInput carries the validated request data for creating a user.
type CreateUserInput struct {
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
	Role      models.Roles
}

// CreateUser enforces unique username/email, hashes the password, and persists
// the new account.
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	if !isValidRole(input.Role) {
		return nil, fmt.Errorf("%w: invalid role %q", models.ErrBadRequest, input.Role)
	}

	if _, err := s.userRepo.GetByUsername(ctx, s.db, input.Username); err == nil {
		return nil, fmt.Errorf("%w: username already exists", models.ErrConflict)
	} else if !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	if _, err := s.userRepo.GetByEmail(ctx, s.db, input.Email); err == nil {
		return nil, fmt.Errorf("%w: email already exists", models.ErrConflict)
	} else if !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	hash, err := authutils.HashPassword(input.Password)
	if err != nil {
		logger.LogError(err, "failed to hash password during user creation", logger.Fields{
			"username": input.Username,
		})
		return nil, err
	}

	user := &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		PassHash:  hash,
		Role:      input.Role,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, s.db, user); err != nil {
		logger.LogError(err, "failed to persist new user", logger.Fields{
			"username": input.Username,
		})
		return nil, err
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("role", string(user.Role)).
		Msg("user created")

	return user, nil
}

// GetUserByID returns a single user by primary key.
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, s.db, id)
}

// ListUsers returns one page of users and the total count.
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.List(ctx, s.db, page, pageSize)
}

// UpdateUserRole changes a user's assigned role after validating it is known.
func (s *UserService) UpdateUserRole(ctx context.Context, id uuid.UUID, role models.Roles) (*models.User, error) {
	if !isValidRole(role) {
		return nil, fmt.Errorf("%w: invalid role %q", models.ErrBadRequest, role)
	}

	user, err := s.userRepo.GetByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}

	user.Role = role
	if err := s.userRepo.Update(ctx, s.db, user); err != nil {
		logger.LogError(err, "failed to update user role", logger.Fields{"user_id": id.String()})
		return nil, err
	}

	log.Info().
		Str("user_id", id.String()).
		Str("role", string(role)).
		Msg("user role updated")

	return user, nil
}

func isValidRole(role models.Roles) bool {
	switch role {
	case models.Admin, models.Manager, models.Staff:
		return true
	default:
		return false
	}
}
