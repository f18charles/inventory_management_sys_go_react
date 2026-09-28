package user

import (
	"context"
	"errors"
	"fmt"

	"i_m_s/internal/models"
	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Service owns user lifecycle business rules: uniqueness, role validity, and
// password hashing. It depends on the domain's Repository interface.
type Service struct {
	db   *gorm.DB
	repo Repository
}

func NewService(db *gorm.DB, repo Repository) *Service {
	return &Service{db: db, repo: repo}
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
func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	if !isValidRole(input.Role) {
		return nil, fmt.Errorf("%w: invalid role %q", models.ErrBadRequest, input.Role)
	}

	if _, err := s.repo.GetByUsername(ctx, s.db, input.Username); err == nil {
		return nil, fmt.Errorf("%w: username already exists", models.ErrConflict)
	} else if !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	if _, err := s.repo.GetByEmail(ctx, s.db, input.Email); err == nil {
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

	created := &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		PassHash:  hash,
		Role:      input.Role,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, s.db, created); err != nil {
		logger.LogError(err, "failed to persist new user", logger.Fields{
			"username": input.Username,
		})
		return nil, err
	}

	log.Info().
		Str("user_id", created.ID.String()).
		Str("role", string(created.Role)).
		Msg("user created")

	return created, nil
}

// GetByID returns a single user by primary key.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, s.db, id)
}

// List returns one page of users and the total count.
func (s *Service) List(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return s.repo.List(ctx, s.db, page, pageSize)
}

// UpdateRole changes a user's assigned role after validating it is known.
func (s *Service) UpdateRole(ctx context.Context, id uuid.UUID, role models.Roles) (*models.User, error) {
	if !isValidRole(role) {
		return nil, fmt.Errorf("%w: invalid role %q", models.ErrBadRequest, role)
	}

	found, err := s.repo.GetByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}

	found.Role = role
	if err := s.repo.Update(ctx, s.db, found); err != nil {
		logger.LogError(err, "failed to update user role", logger.Fields{"user_id": id.String()})
		return nil, err
	}

	log.Info().
		Str("user_id", id.String()).
		Str("role", string(role)).
		Msg("user role updated")

	return found, nil
}

func isValidRole(role models.Roles) bool {
	switch role {
	case models.Admin, models.Manager, models.Staff:
		return true
	default:
		return false
	}
}
