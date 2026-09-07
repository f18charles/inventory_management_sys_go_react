package config_test

import (
	"os"
	"testing"

	"i_m_s/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	// Clear any existing env variables that might override defaults
	envKeys := []string{
		"APP_ENV", "APP_PORT", "DB_HOST", "DB_PORT",
		"DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"JWT_SECRET", "JWT_EXPIRATION",
	}
	for _, key := range envKeys {
		os.Unsetenv(key)
	}

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "inventory", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
}

func TestConfigDSN(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "127.0.0.1",
		DBPort:     "5433",
		DBUser:     "testuser",
		DBPassword: "secretpassword",
		DBName:     "testdb",
		DBSSLMode:  "require",
	}

	expectedDSN := "host=127.0.0.1 port=5433 user=testuser password=secretpassword dbname=testdb sslmode=require"
	assert.Equal(t, expectedDSN, cfg.DSN())
}
