package utils_test

import (
	"errors"
	"testing"

	"i_m_s/internal/utils/logger"

	"github.com/stretchr/testify/assert"
)

func TestLoggerInitAndLogError(t *testing.T) {
	// Initialize logger in development mode
	logger.Init("development")

	// Call LogError with non-sensitive and sensitive fields
	assert.NotPanics(t, func() {
		logger.LogError(errors.New("test error"), "test failure message", logger.Fields{
			"user_id":  "12345",
			"password": "super-secret-password", // should be redacted
		})
	})

	// Initialize logger in production mode
	logger.Init("production")

	assert.NotPanics(t, func() {
		logger.LogError(errors.New("test error prod"), "test safe failure message", logger.Fields{
			"operation": "create_user",
			"token":     "jwt-token-string", // should be redacted
		})
	})
}
