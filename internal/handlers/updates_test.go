package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/middlewares"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
)

var loggerOnce sync.Once

func TestMain(m *testing.M) {
	loggerOnce.Do(func() {
		logger.Initialize("info")
	})

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestUpdates(t *testing.T) {
	type args struct {
		//w       http.ResponseWriter
		r       *http.Request
		storage storage.Storage
	}
	tests := []struct {
		name            string
		args            args
		metricName      string
		expectedStatus  int
		expectedMetrics map[string]utils.Metrics
	}{
		{
			name: "Positive #1",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(
					`[
				{
					"id": "RandomValueCounter",
					"type": "counter",
					"delta": 2843918916068879
				},
				{
					"id": "RandomValueGauge",
					"type": "gauge",
					"value": 0.2843918916068879
				}
			]`)),
				storage: storage.NewMemStorage(&sync.Map{}),
			},
			expectedStatus: http.StatusOK,
			metricName:     "testCounter",
			expectedMetrics: map[string]utils.Metrics{
				"RandomValueCounter": utils.NewMetrics("RandomValueCounter", 2843918916068879, true),
				"RandomValueGauge":   utils.NewMetrics("RandomValueGauge", 0.2843918916068879, false),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rr := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = tt.args.r
			r := gin.Default()

			r.RedirectTrailingSlash = false
			r.POST("/updates", middlewares.WithLogging(), func(ctx *gin.Context) {
				Updates(ctx, tt.args.storage)
			})
			r.HandleContext(ctx)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				metrics := tt.args.storage.GetAllMetrics()
				assert.Equal(t, tt.expectedMetrics, metrics)
			}
		})
	}
}
