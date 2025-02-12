package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)


func WithLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		


		c.Next()

		duration := time.Since(start)

		logger.Log.Info("Request received",
			zap.String("URI", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Duration("duration", duration),
		)

		logger.Log.Info("Response sent",
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}