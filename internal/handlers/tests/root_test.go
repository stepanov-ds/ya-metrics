package handlers

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/handlers"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name           string
		st             storage.Storage
		expectedStatus int
		expectedBody   string
		fillStorage    bool
		rr             *http.Request
	}{
		// TODO: Add test cases.
		{
			name:           "Positive #1 Get empty storage",
			st:             storage.NewMemStorage(&sync.Map{}),
			expectedStatus: http.StatusOK,
			expectedBody:   "\"{}\"",
			fillStorage:    false,
		},
		{
			name:           "Positive #2 Get storage with rewrited Counter and rewrited Gauge",
			st:             storage.NewMemStorage(&sync.Map{}),
			expectedStatus: http.StatusOK,
			expectedBody:   "\"{\\\"test1\\\":{\\\"delta\\\":5},\\\"test2\\\":{\\\"value\\\":2.2}}\"",
			fillStorage:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fillStorage {
				tt.st.SetMetricCounter("test1", 1)
				tt.st.SetMetricCounter("test1", 4)
				tt.st.SetMetricGauge("test2", 1.1)
				tt.st.SetMetricGauge("test2", 2.2)
			}

			gin.SetMode(gin.TestMode)
			rr := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			r := gin.Default()

			r.RedirectTrailingSlash = false
			r.GET("/", func(c *gin.Context) {
				handlers.Root(c, tt.st)
			})
			r.HandleContext(ctx)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
