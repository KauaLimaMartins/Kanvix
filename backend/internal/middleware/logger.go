package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()
		reqID, _ := c.Get(RequestIDHeader)

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		log.Info("request",
			"request_id", reqID,
			"status", status,
			"method", method,
			"path", path,
			"latency_ms", latency.Milliseconds(),
			"size", size,
		)
	}
}

