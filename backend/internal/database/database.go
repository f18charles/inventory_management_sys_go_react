package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Connect(dsn string, isDev bool) (*gorm.DB, error) {
	logLevel := gormlogger.Silent
	if isDev {
		logLevel = gormlogger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
		// Translate dialect errors (duplicate key, FK, check constraint) into
		// GORM sentinel errors so repositories can map them without string
		// matching or leaking driver-specific types upward.
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Connection pool tuning
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
