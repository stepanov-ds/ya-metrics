package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/cmd/server/handlers"
	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
	"github.com/stepanov-ds/ya-metrics/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	type args struct {
		//w       http.ResponseWriter
		r       *http.Request
		storage storage.Storage
	}
	tests := []struct {
		name           string
		args           args
		metricName     string
		expectedStatus int
		expectedMetric utils.Metric
	}{
		// TODO: Add test cases.
		{
			name: "Negative #1 Method GET not alowed",
			args: args{
				r:       httptest.NewRequest(http.MethodGet, "/update/gauge/testGauge/23.1", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "Negative #2 no metric value for counter",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #3 no metric value for gauge",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #4 wrong value for counter (float)",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/123.3", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative #5 wrong value for counter (string)",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative #6 wrong value for gauge (string)",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative #7 no metric name for gauge",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/gauge//asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #8 no metric name for counter",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/counter//asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #9 wrong metric type",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/asd/testAsd/123", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Positive #1 counter",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/123", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusOK,
			metricName:     "testCounter",
			expectedMetric: utils.Metric{
				Counter:   123,
				Gauge:     0,
				IsCounter: true,
			},
		},
		{
			name: "Positive #2 gauge",
			args: args{
				r:       httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/123.1/", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusOK,
			metricName:     "testGauge",
			expectedMetric: utils.Metric{
				Counter:   0,
				Gauge:     123.1,
				IsCounter: false,
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
			r.Any("/update/:metric_type/:metric_name/:value/", func(c *gin.Context) {
				handlers.Update(c, tt.args.storage)
			})
			r.Any("/update/:metric_type/:metric_name/:value", func(c *gin.Context) {
				handlers.Update(c, tt.args.storage)
			})
			r.HandleContext(ctx)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				metric, found := tt.args.storage.GetMetric(tt.metricName)
				if found {
					assert.Equal(t, tt.expectedMetric, metric)
				} else {
					assert.Fail(t, "metric not found in storage")
				}
			}
		})
	}
}
