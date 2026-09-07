// internal/utils/logger/logger.go
//
// One switch (APP_ENV) controls two logging behaviors at once:
//   - FORMAT: pretty/colorized console in dev, plain JSON in production
//     (JSON is what log aggregators expect; console is what a human wants
//     while actively developing).
//   - VERBOSITY: production logs the error's safe, high-level message only;
//     development additionally attaches the raw underlying error detail
//     and a stack trace. This is the actual answer to "hide sensitive
//     detail in prod, but let me flip a switch to see the exact issue":
//     restart the service with APP_ENV=development (e.g. against a
//     staging replica, or briefly in prod while actively debugging) and
//     the exact same log call now includes the full detail.
//
// Regardless of mode, known-sensitive keys are ALWAYS stripped from
// structured fields — the dev/prod switch controls verbosity of the error
// detail itself, not an excuse to ever log a password or token.

package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var isDev bool

// sensitiveKeys are never logged, in either mode. Extend this list rather
// than working around it.
var sensitiveKeys = map[string]struct{}{
	"password":      {},
	"pass_hash":     {},
	"password_hash": {},
	"token":         {},
	"authorization": {},
	"jwt_secret":    {},
	"secret":        {},
}

// Init sets the global zerolog logger based on APP_ENV. Call this once at
// startup, before anything else logs.
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

// LogError is the standard way to log a failure from a service or
// repository. `safe` is a short, human-written summary safe to appear in
// any environment (e.g. "sale transaction failed"). `err` is the actual
// Go error — its full text is only attached in development.
func LogError(err error, safe string, fields map[string]any) {
	event := log.Error().Err(errIf(isDev, err))
	for k, v := range fields {
		if _, blocked := sensitiveKeys[k]; blocked {
			continue
		}
		event = event.Interface(k, v)
	}
	if !isDev {
		// In production, still record that a real error occurred and its
		// safe classification, without ever writing err.Error() itself —
		// the raw text may contain SQL fragments, file paths, or values
		// pulled from the request.
		event = event.Bool("detail_suppressed", true)
	}
	event.Msg(safe)
}

// errIf returns the error only in dev; nil in production so zerolog's
// Err() call attaches nothing.
func errIf(dev bool, err error) error {
	if dev {
		return err
	}
	return nil
}

// Fields is a small convenience constructor so call sites read cleanly:
// logger.LogError(err, "sale transaction failed, rolled back", logger.Fields{
//     "customer_id": input.CustomerID,
// })
type Fields = map[string]any
