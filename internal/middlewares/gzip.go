// Package middlewares implements custom middleware functions for the Gin router.
//
// Currently provides Gzip middleware for handling compressed request bodies.
package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		c.Next()
	}
}
