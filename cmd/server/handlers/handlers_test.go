package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stepanov-ds/ya-metrics/cmd/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	type args struct {
		w       http.ResponseWriter
		r       *http.Request
		storage storage.Repositories
	}
	tests := []struct {
		name string
		args args
		metricName string
		expectedStatus int
		expectedMetric storage.Metric
	}{
		// TODO: Add test cases.
		{
			name: "Negative #1 Method GET not alowed",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/update/gauge/testGauge/23.1", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "Negative #2 no metric value for counter",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #3 no metric value for gauge",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #4 wrong value for counter (float)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/123.3", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative #5 wrong value for counter (string)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative #6 wrong value for gauge (string)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative #7 no metric name for gauge",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge//asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #8 no metric name for counter",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter//asd", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative #9 wrong metric type",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/asd/testAsd/123", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Positive #1 counter",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/123", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusOK,
			metricName: "testCounter",
			expectedMetric: storage.Metric{
				Counter: 123,
				Gauge: 0,
				IsCounter: true,
			},
		},
		{
			name: "Positive #2 gauge",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/123.1", nil),
				storage: storage.NewMemStorage(),
			},
			expectedStatus: http.StatusOK,
			metricName: "testGauge",
			expectedMetric: storage.Metric{
				Counter: 0,
				Gauge: 123.1,
				IsCounter: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Update(tt.args.w, tt.args.r, tt.args.storage)

			rr := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if (tt.expectedStatus == http.StatusOK) {
				assert.Equal(t, tt.expectedMetric, tt.args.storage.GetMetric(tt.metricName))
			}
		})
	}
}
