package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	gzipWriter *gzip.Writer
	buffer     *bytes.Buffer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	// Записываем данные в gzipWriter
	return w.gzipWriter.Write(data)
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "gzip" {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}

			reader, err := gzip.NewReader(bytes.NewReader(body))
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			defer reader.Close()

			decompressed, err := io.ReadAll(reader)
			if err != nil {
				print(err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(decompressed))
		}

		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)

			writer := c.Writer
			writer.Header().Set("Content-Encoding", "gzip")

			c.Writer = &gzipResponseWriter{
				ResponseWriter: writer,
				gzipWriter:     gz,
				buffer:         &buf,
			}

			defer func() {
				gz.Close()
				writer.Write(buf.Bytes())
			}()
		}
		c.Next()
	}
}
