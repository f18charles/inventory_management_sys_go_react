package utils_test

import (
	"strings"
	"testing"
	"time"

	"i_m_s/internal/models"
	authutils "i_m_s/internal/utils/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPasswordAndCheck(t *testing.T) {
	hash, err := authutils.HashPassword("s3cret-password")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "s3cret-password", hash)
	assert.True(t, authutils.CheckPassword(hash, "s3cret-password"))
	assert.False(t, authutils.CheckPassword(hash, "wrong-password"))
}

func TestHashPasswordIsSalted(t *testing.T) {
	h1, err := authutils.HashPassword("same-password")
	require.NoError(t, err)
	h2, err := authutils.HashPassword("same-password")
	require.NoError(t, err)
	assert.NotEqual(t, h1, h2)
}

func TestJWTGenerateAndParse(t *testing.T) {
	mgr := authutils.NewJWTManager("test-secret", time.Hour)
	userID := uuid.New()

	token, err := mgr.Generate(userID, models.Manager)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := mgr.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, models.Manager, claims.Role)
}

func TestJWTParseRejectsExpired(t *testing.T) {
	mgr := authutils.NewJWTManager("test-secret", -time.Minute)
	token, err := mgr.Generate(uuid.New(), models.Staff)
	require.NoError(t, err)

	_, err = mgr.Parse(token)
	assert.ErrorIs(t, err, authutils.ErrInvalidToken)
}

func TestJWTParseRejectsWrongSecret(t *testing.T) {
	issuer := authutils.NewJWTManager("secret-a", time.Hour)
	verifier := authutils.NewJWTManager("secret-b", time.Hour)

	token, err := issuer.Generate(uuid.New(), models.Admin)
	require.NoError(t, err)

	_, err = verifier.Parse(token)
	assert.ErrorIs(t, err, authutils.ErrInvalidToken)
}

func TestJWTParseRejectsMalformed(t *testing.T) {
	mgr := authutils.NewJWTManager("test-secret", time.Hour)

	for _, bad := range []string{"", "not-a-token", "a.b.c", strings.Repeat("x", 16)} {
		_, err := mgr.Parse(bad)
		assert.ErrorIs(t, err, authutils.ErrInvalidToken)
	}
}
