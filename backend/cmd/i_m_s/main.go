package main

import (
	"fmt"

	"i_m_s/internal/config"
	"i_m_s/internal/database"
	"i_m_s/internal/handlers"
	"i_m_s/internal/utils/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Initialize logger (must happen before any logging calls)
	logger.Init(cfg.AppEnv)

	log.Info().
		Str("app_env", cfg.AppEnv).
		Str("app_port", cfg.AppPort).
		Msg("starting inventory management system")

	// Connect to PostgreSQL
	db, err := database.Connect(cfg.DSN(), cfg.AppEnv != "production")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	log.Info().Msg("database connection established")

	// Setup router with middleware
	router := handlers.SetupRouter(db, cfg.AppEnv != "production")

	// Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Info().Str("address", addr).Msg("HTTP server listening")
	if err := router.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
