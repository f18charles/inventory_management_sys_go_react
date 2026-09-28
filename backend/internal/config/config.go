package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	// DatabaseURL is the application database connection string (Neon
	// PostgreSQL). When set it takes precedence over the discrete DB_* values.
	DatabaseURL string

	// TestDatabaseURL is the local PostgreSQL connection string used by tests.
	// Integration tests use this instead of the application database.
	TestDatabaseURL string

	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiration string
}

func Load() (*Config, error) {
	// Attempt to load .env file; ignore if missing
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	cfg := &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		AppPort:         getEnv("APP_PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		TestDatabaseURL: getEnv("TEST_DATABASE_URL", ""),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "inventory"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "default-dev-secret-change-in-prod"),
		JWTExpiration:   getEnv("JWT_EXPIRATION", "24h"),
	}

	return cfg, nil
}

// DSN returns the application database connection string: the Neon
// DATABASE_URL when provided, otherwise a DSN built from the discrete DB_*
// values.
func (c *Config) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return c.buildDSN()
}

// TestDSN returns the connection string tests run against: the local
// TEST_DATABASE_URL when provided, otherwise the discrete DB_* values.
func (c *Config) TestDSN() string {
	if c.TestDatabaseURL != "" {
		return c.TestDatabaseURL
	}
	return c.buildDSN()
}

func (c *Config) buildDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// JWTExpiry parses JWT_EXPIRATION, falling back to 24h when unset or invalid.
func (c *Config) JWTExpiry() time.Duration {
	expiry, err := time.ParseDuration(c.JWTExpiration)
	if err != nil || expiry <= 0 {
		return 24 * time.Hour
	}
	return expiry
}

func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return fallback
}
