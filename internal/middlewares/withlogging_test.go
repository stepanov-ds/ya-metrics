package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	logger.Initialize("fatal")
	oldLogger := logger.Log
	defer func() { logger.Log = oldLogger }()

	return r
}

func TestWithLogging_Middleware(t *testing.T) {
	r := setupTestRouter()

	logger.Log = zap.NewNop()
	defer func() {
		logger.Log = nil
	}()

	r.Use(WithLogging())

	r.POST("/log", func(c *gin.Context) {
		c.String(http.StatusOK, "Response body")
	})

	reqBody := strings.NewReader(`{"key":"value"}`)
	req, _ := http.NewRequest("POST", "/log", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Response body", w.Body.String())
}
