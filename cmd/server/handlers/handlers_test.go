package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stepanov-ds/ya-metrics/cmd/server/metric"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	type args struct {
		w       http.ResponseWriter
		r       *http.Request
		storage *metric.MemStorage
	}
	tests := []struct {
		name string
		args args
		expectedStatus int
		expectedStorage *metric.MemStorage
	}{
		// TODO: Add test cases.
		{
			name: "Negative #1 Method GET not alowed",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/update/gauge/testGauge/23.1", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #2 no metric value for counter",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #3 no metric value for gauge",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #4 wrong value for counter (float)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/123.3", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #5 wrong value for counter (string)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/asd", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #6 wrong value for gauge (string)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/asd", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #7 no metric name for gauge",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge//asd", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #8 no metric name for counter",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter//asd", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusNotFound,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Negative #9 wrong metric type",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/asd/testAsd/123", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedStorage: metric.NewMemStorage(),
		},
		{
			name: "Positive #1 counter",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/123", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusOK,
			expectedStorage: &metric.MemStorage{
				Storage: map[string]metric.Metric{
					"testCounter": {
						Counter: 123,
						Gauge: 0,
						IsCounter: true,
					},
				},
			},
		},
		{
			name: "Positive #2 gauge",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/123.1", nil),
				storage: metric.NewMemStorage(),
			},
			expectedStatus: http.StatusOK,
			expectedStorage: &metric.MemStorage{
				Storage: map[string]metric.Metric{
					"testGauge": {
						Counter: 0,
						Gauge: 123.1,
						IsCounter: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Update(tt.args.w, tt.args.r, tt.args.storage)

			rr := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedStorage, tt.args.storage)
		})
	}
}
