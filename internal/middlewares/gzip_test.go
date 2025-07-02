package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func compressWithGzip(input string) io.Reader {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(input))
	if err != nil {
		panic(err)
	}
	zw.Close()
	return bytes.NewReader(buf.Bytes())
}

func setupTestRouterWithGzip() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Gzip())

	r.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})

	return r
}

func Test_GzipMiddleware_NoCompression(t *testing.T) {
	r := setupTestRouterWithGzip()

	reqBody := strings.NewReader("plain text body")
	req, _ := http.NewRequest("POST", "/test", reqBody)
	req.Header.Set("Content-Type", "text/plain")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "plain text body", resp.Body.String())
}

func Test_GzipMiddleware_ValidCompression(t *testing.T) {
	r := setupTestRouterWithGzip()

	const original = "this is a test string"

	reqBody := compressWithGzip(original)
	req, _ := http.NewRequest("POST", "/test", reqBody)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Encoding", "gzip")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, original, resp.Body.String())
}

func Test_GzipMiddleware_UnsupportedEncoding(t *testing.T) {
	r := setupTestRouterWithGzip()

	reqBody := strings.NewReader("some data")
	req, _ := http.NewRequest("POST", "/test", reqBody)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Encoding", "deflate")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "some data", resp.Body.String())
}
