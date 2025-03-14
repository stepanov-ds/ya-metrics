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

		originalWriter := c.Writer
		var bodyBuf bytes.Buffer
		loggedWriter := &LoggedResponseWriter{
			ResponseWriter: originalWriter,
			Body:           &bodyBuf,
		}
		c.Writer = loggedWriter

		c.Next()

		respBody := bodyBuf.String()

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
			zap.String("body", respBody),
		)
	}
}

type LoggedResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *LoggedResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *LoggedResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}
