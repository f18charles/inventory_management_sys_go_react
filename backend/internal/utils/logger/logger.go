package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var isDev bool

var sensitiveKeys = map[string]struct{}{
	"password":      {},
	"pass_hash":     {},
	"password_hash": {},
	"token":         {},
	"authorization": {},
	"jwt_secret":    {},
	"secret":        {},
}

func Init(env string) {
	isDev = env != "production"

	if isDev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
		log.Info().Msg("logger initialized in DEVELOPMENT mode (verbose, error detail included)")
		return
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Info().Msg("logger initialized in PRODUCTION mode (error detail suppressed)")
}

func LogError(err error, safe string, fields map[string]any) {
	event := log.Error().Err(errIf(isDev, err))
	for k, v := range fields {
		if _, blocked := sensitiveKeys[k]; blocked {
			continue
		}
		event = event.Interface(k, v)
	}
	if !isDev {
		event = event.Bool("detail_suppressed", true)
	}
	event.Msg(safe)
}

func errIf(dev bool, err error) error {
	if dev {
		return err
	}
	return nil
}

type Fields = map[string]any
