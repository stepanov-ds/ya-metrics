package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/handlers"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
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
			st:             storage.NewMemStorage(),
			expectedStatus: http.StatusOK,
			expectedBody:   "\"{}\"",
			fillStorage:    false,
		},
		{
			name:           "Positive #2 Get storage with rewrited Counter and rewrited Gauge",
			st:             storage.NewMemStorage(),
			expectedStatus: http.StatusOK,
			expectedBody:   "\"{\\\"test1\\\":{\\\"Counter\\\":5},\\\"test2\\\":{\\\"Gauge\\\":2.2}}\"",
			fillStorage:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fillStorage {
				tt.st.SetMetric("test1", utils.NewMetricCounter(1))
				tt.st.SetMetric("test1", utils.NewMetricCounter(4))
				tt.st.SetMetric("test2", utils.NewMetricGauge(1.1))
				tt.st.SetMetric("test2", utils.NewMetricGauge(2.2))
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
