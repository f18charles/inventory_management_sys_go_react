package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		event := log.Info()
		if statusCode >= 400 && statusCode < 500 {
			event = log.Warn()
		} else if statusCode >= 500 {
			event = log.Error()
		}

		event.
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Msg("http request completed")
	}
}
