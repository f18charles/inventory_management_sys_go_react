package auth

import (
	"errors"
	"fmt"
	"time"

	"i_m_s/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidToken is returned for any token that is malformed, tampered with,
// expired, or signed with the wrong key. Callers map it to HTTP 401.
var ErrInvalidToken = errors.New("invalid token")

const tokenIssuer = "i_m_s"

// Claims is the custom JWT payload: the authenticated user's identity and role
// on top of the standard registered claims.
type Claims struct {
	UserID string       `json:"user_id"`
	Role   models.Roles `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager issues and validates HS256 tokens using a shared secret.
type JWTManager struct {
	secret     []byte
	expiration time.Duration
}

// NewJWTManager builds a manager. expiration is taken from configuration
// (JWT_EXPIRATION) so token lifetime is tunable per environment.
func NewJWTManager(secret string, expiration time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), expiration: expiration}
}

// Generate signs a token for the given user and role.
func (m *JWTManager) Generate(userID uuid.UUID, role models.Roles) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// Parse validates a token and returns its claims. All failures collapse to
// ErrInvalidToken so callers never branch on jwt library internals.
func (m *JWTManager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
