package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	tests := []struct {
		st             storage.Storage
		request        *http.Request
		name           string
		expectedBody   string
		expectedStatus int
		fillStorage    bool
	}{
		// TODO: Add test cases.
		{
			name:           "Negative #1 metric not exist",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			request:        httptest.NewRequest(http.MethodGet, "/value/gauge/test", nil),
		},
		{
			name:           "Negative #2 no metric name",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			request:        httptest.NewRequest(http.MethodGet, "/value/gauge/", nil),
		},
		{
			name:           "Negative #3 no metric name 2",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			request:        httptest.NewRequest(http.MethodGet, "/value/gauge", nil),
		},
		{
			name:           "Negative #4 no metric name 3",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			request:        httptest.NewRequest(http.MethodGet, "/value/gauge//", nil),
		},
		{
			name:           "Negative #5 wrong type",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			request:        httptest.NewRequest(http.MethodGet, "/value/gage/test", nil),
		},
		{
			name:           "Negative #6 try to get Counter rewrited by Gauge",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			request:        httptest.NewRequest(http.MethodGet, "/value/counter/test3", nil),
		},
		{
			name:           "Negative #7 try to get Gauge rewrited by Counter",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			request:        httptest.NewRequest(http.MethodGet, "/value/gauge/test4", nil),
		},
		{
			name:           "Positive #1 get rewrited Counter by another Counter",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusOK,
			expectedBody:   "5",
			request:        httptest.NewRequest(http.MethodGet, "/value/counter/test1", nil),
		},
		{
			name:           "Positive #2 get rewrited Gauge by another Gauge",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusOK,
			expectedBody:   "2.2",
			request:        httptest.NewRequest(http.MethodGet, "/value/gauge/test2", nil),
		},
		{
			name:           "Positive #3 get rewrited Counter by Gauge",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusOK,
			expectedBody:   "6.6",
			request:        httptest.NewRequest(http.MethodGet, "/value/gauge/test3", nil),
		},
		{
			name:           "Positive #4 get rewrited Gauge by Counter",
			st:             storage.NewMemStorage(&sync.Map{}),
			fillStorage:    true,
			expectedStatus: http.StatusOK,
			expectedBody:   "7",
			request:        httptest.NewRequest(http.MethodGet, "/value/counter/test4", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fillStorage {
				tt.st.SetMetric(context.Background(), "test1", 1, true)
				tt.st.SetMetric(context.Background(), "test1", 4, true)
				tt.st.SetMetric(context.Background(), "test2", 1.1, false)
				tt.st.SetMetric(context.Background(), "test2", 2.2, false)
				tt.st.SetMetric(context.Background(), "test3", 6, true)
				tt.st.SetMetric(context.Background(), "test3", 6.6, false)
				tt.st.SetMetric(context.Background(), "test4", 7.7, false)
				tt.st.SetMetric(context.Background(), "test4", 7, true)
			}
			gin.SetMode(gin.TestMode)
			rr := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = tt.request
			r := gin.Default()

			r.RedirectTrailingSlash = false
			r.GET("/value/:metric_type/:metric_name", func(ctx *gin.Context) {
				Value(ctx, tt.st)
			})
			r.GET("/value/:metric_type/:metric_name/", func(ctx *gin.Context) {
				Value(ctx, tt.st)
			})
			r.HandleContext(ctx)
			// print("_______________________________________BODY___________________")
			// print(rr.Body.String())
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
