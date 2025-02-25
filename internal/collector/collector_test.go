package collector

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCollector_CollectMetrics(t *testing.T) {
	tests := []struct {
		name      string
		c         *Collector
		pollCount int64
	}{
		// TODO: Add test cases.
		{
			name:      "Positive #1 collect once",
			c:         NewCollector(&sync.Map{}),
			pollCount: 1,
		},
		{
			name:      "Positive #2 collect twice",
			c:         NewCollector(&sync.Map{}),
			pollCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < int(tt.pollCount); i++ {
				tt.c.CollectMetrics()
			}
			if value, ok := tt.c.Metrics.Load("PollCount"); ok {
				if v, ok := value.(utils.Metrics); ok {
					assert.Equal(t, tt.pollCount, v.Get())
				} else {
					assert.Fail(t, "metric PollCount is not utils.MetricCounter, metric PollCounter is ", reflect.TypeOf(value))
				}
			} else {
				assert.Fail(t, "metric PollCount does not exist")
			}
		})
	}
}
