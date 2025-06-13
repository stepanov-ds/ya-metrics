// Package middlewares implements custom middleware functions for the Gin router.
//
// Currently provides Gzip middleware for handling compressed request bodies.
package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	// "strings"
	// "sync"

	"github.com/gin-gonic/gin"
)


// var gzipPool = sync.Pool{
// 	New: func() any {
// 		// var buf bytes.Buffer
// 		w := gzip.NewWriter(io.Discard)
// 		return w
// 	},
// }

// type gzipResponseWriter struct {
// 	gin.ResponseWriter
// 	gzipWriter *gzip.Writer
// 	// buffer     *bytes.Buffer
// }

// func (w *gzipResponseWriter) Write(data []byte) (int, error) {
// 	return w.gzipWriter.Write(data)
// }

// Gzip returns a Gin middleware handler that manages gzip-compressed requests.
//
// It automatically decompresses incoming gzip-encoded payloads so handlers can
// read them as plain text. Useful for APIs that accept compressed data.
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

		// if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
		// 	// var buf bytes.Buffer
		// 	// gz := gzip.NewWriter(&buf)

		// 	gz := gzipPool.Get().(*gzip.Writer)
		// 	defer gzipPool.Put(gz)
		// 	gz.Reset(c.Writer)

		// 	c.Writer.Header().Set("Content-Encoding", "gzip")
		// 	c.Writer.Header().Set("Vary", "Accept-Encoding")

		// 	c.Writer = &gzipResponseWriter{
		// 		ResponseWriter: c.Writer,
		// 		gzipWriter:     gz,
		// 		// buffer:         &buf,
		// 	}

		// 	defer func() {
		// 		gz.Close()
		// 		// writer.Write(buf.Bytes())
		// 	}()
		// }
		c.Next()
	}
}
