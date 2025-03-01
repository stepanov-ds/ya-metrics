package middlewares

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)

func WithLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to read body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()

		duration := time.Since(start)

		logger.Log.Info("Request received",
			zap.String("URI", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Duration("duration", duration),
			zap.String("body", string(body)),
		)

		logger.Log.Info("Response sent",
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}
