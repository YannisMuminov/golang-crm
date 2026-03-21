package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		status := c.Writer.Status()

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
		}

		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		switch {
		case status >= 500:
			L.Error("request", fields...)
		case status >= 400:
			L.Warn("request", fields...)
		default:
			L.Info("request", fields...)
		}
	}
}
