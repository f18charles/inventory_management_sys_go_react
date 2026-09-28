package auth

import (
	"context"
	"errors"

	"i_m_s/internal/models"
	"i_m_s/internal/user"
	authutils "i_m_s/internal/utils/auth"
	"i_m_s/internal/utils/logger"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Service verifies credentials and issues JWT tokens. It depends on the user
// domain's Repository interface (auth -> user in the domain DAG) and never
// touches the Gin context or builds HTTP status codes.
type Service struct {
	db    *gorm.DB
	users user.Repository
	jwt   *authutils.JWTManager
}

func NewService(db *gorm.DB, users user.Repository, jwt *authutils.JWTManager) *Service {
	return &Service{db: db, users: users, jwt: jwt}
}

// Login authenticates a user by username or email and returns the account plus
// a signed token. Unknown users and wrong passwords deliberately collapse to
// the same error to avoid account enumeration.
func (s *Service) Login(ctx context.Context, usernameOrEmail, password string) (*models.User, string, error) {
	account, err := s.users.GetByUsername(ctx, s.db, usernameOrEmail)
	if errors.Is(err, models.ErrNotFound) {
		account, err = s.users.GetByEmail(ctx, s.db, usernameOrEmail)
	}
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			log.Warn().Str("identifier", usernameOrEmail).Msg("login failed: unknown user")
			return nil, "", models.ErrInvalidCredentials
		}
		logger.LogError(err, "login lookup failed", logger.Fields{"identifier": usernameOrEmail})
		return nil, "", err
	}

	if !account.IsActive {
		log.Warn().Str("user_id", account.ID.String()).Msg("login failed: account inactive")
		return nil, "", models.ErrAccountInactive
	}

	if !authutils.CheckPassword(account.PassHash, password) {
		log.Warn().Str("user_id", account.ID.String()).Msg("login failed: invalid password")
		return nil, "", models.ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(account.ID, account.Role)
	if err != nil {
		logger.LogError(err, "failed to issue token", logger.Fields{"user_id": account.ID.String()})
		return nil, "", err
	}

	log.Info().
		Str("user_id", account.ID.String()).
		Str("role", string(account.Role)).
		Msg("login successful")

	return account, token, nil
}

// GetProfile returns the authenticated user's account.
func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.users.GetByID(ctx, s.db, userID)
}
