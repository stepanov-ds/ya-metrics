package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
)

func setupTestRouterWithHashCheck(key string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	server.Key = &key

	logger.Initialize("fatal")

	r := gin.New()
	r.Use(HashCheck())

	r.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, "Response: "+string(body))
	})

	return r
}

func Test_HashCheck_NoKey(t *testing.T) {
	r := setupTestRouterWithHashCheck("")

	reqBody := strings.NewReader(`{"value": 3.14}`)
	req, _ := http.NewRequest("POST", "/test", reqBody)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Header().Get("HashSHA256"))
}

func Test_HashCheck_InvalidHash(t *testing.T) {
	key := "secret_key"
	r := setupTestRouterWithHashCheck(key)

	body := `{"value": 3.14}`
	validHash := utils.CalculateHashWithKey([]byte(body), key)

	wrongHash := validHash[:len(validHash)-1] + "a"

	reqBody := strings.NewReader(body)
	req, _ := http.NewRequest("POST", "/test", reqBody)
	req.Header.Set("HashSHA256", wrongHash)

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func Test_HashCheck_ValidHash(t *testing.T) {
	key := "secret_key"
	r := setupTestRouterWithHashCheck(key)

	body := `{"value": 3.14}`
	hash := utils.CalculateHashWithKey([]byte(body), key)

	reqBody := strings.NewReader(body)
	req, _ := http.NewRequest("POST", "/test", reqBody)
	req.Header.Set("HashSHA256", hash)

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.NotEmpty(t, resp.Header().Get("HashSHA256"))
	assert.Equal(t, "Response: {\"value\": 3.14}", resp.Body.String())
}

func Test_HashCheck_ResponseHash(t *testing.T) {
	key := "secret_key"
	r := setupTestRouterWithHashCheck(key)

	body := `{"value": 3.14}`
	hash := utils.CalculateHashWithKey([]byte(body), key)

	reqBody := strings.NewReader(body)
	req, _ := http.NewRequest("POST", "/test", reqBody)
	req.Header.Set("HashSHA256", hash)

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	responseBody := "Response: {\"value\": 3.14}"
	calculated := utils.CalculateHashWithKey([]byte(responseBody), key)
	headerHash := resp.Header().Get("HashSHA256")

	assert.Equal(t, calculated, headerHash)
}
