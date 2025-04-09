package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"go.uber.org/zap"
)

func HashCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
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


		hashString := c.GetHeader("HashSHA256")
		
		if utils.CalculateHashWithKey(body, *server.Key) != hashString {
			logger.Log.Error("HashCheck", zap.String("error", "body hash does not match"),
										  zap.String("hashString", hashString),
										  zap.String("calculatedHashString", utils.CalculateHashWithKey(body, *server.Key)))
			c.AbortWithStatusJSON(http.StatusBadRequest, nil)
			return
		}
		c.Next()

		respBody := bodyBuf.String()

		hashRespString := utils.CalculateHashWithKey([]byte(respBody), *server.Key)

		c.Writer.Header().Set("HashSHA256", hashRespString)

	}
}
