package services

import (
	"context"
	"errors"

	"i_m_s/internal/models"
	"i_m_s/internal/repositories"
	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AuthService verifies credentials and issues JWT tokens. It never touches the
// Gin context or builds HTTP status codes.
type AuthService struct {
	db       *gorm.DB
	userRepo repositories.UserRepository
	jwt      *authutils.JWTManager
}

func NewAuthService(db *gorm.DB, userRepo repositories.UserRepository, jwt *authutils.JWTManager) *AuthService {
	return &AuthService{db: db, userRepo: userRepo, jwt: jwt}
}

// Login authenticates a user by username or email and returns the account plus
// a signed token. Unknown users and wrong passwords deliberately collapse to
// the same error to avoid account enumeration.
func (s *AuthService) Login(ctx context.Context, usernameOrEmail, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetByUsername(ctx, s.db, usernameOrEmail)
	if errors.Is(err, models.ErrNotFound) {
		user, err = s.userRepo.GetByEmail(ctx, s.db, usernameOrEmail)
	}
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			log.Warn().Str("identifier", usernameOrEmail).Msg("login failed: unknown user")
			return nil, "", models.ErrInvalidCredentials
		}
		logger.LogError(err, "login lookup failed", logger.Fields{"identifier": usernameOrEmail})
		return nil, "", err
	}

	if !user.IsActive {
		log.Warn().Str("user_id", user.ID.String()).Msg("login failed: account inactive")
		return nil, "", models.ErrAccountInactive
	}

	if !authutils.CheckPassword(user.PassHash, password) {
		log.Warn().Str("user_id", user.ID.String()).Msg("login failed: invalid password")
		return nil, "", models.ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		logger.LogError(err, "failed to issue token", logger.Fields{"user_id": user.ID.String()})
		return nil, "", err
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("role", string(user.Role)).
		Msg("login successful")

	return user, token, nil
}

// GetProfile returns the authenticated user's account.
func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, s.db, userID)
}
