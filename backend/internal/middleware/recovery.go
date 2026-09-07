package middleware

import (
	"fmt"
	"net/http"

	"i_m_s/internal/utils/logger"
	"i_m_s/internal/utils/response"

	"github.com/gin-gonic/gin"
)

func PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				reqID, _ := c.Get("request_id")
				logger.LogError(err, "unexpected panic recovered in middleware", logger.Fields{
					"request_id": reqID,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				})

				response.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
				c.Abort()
			}
		}()

		c.Next()
	}
}
