package collector

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollector_CollectMetrics(t *testing.T) {
	tests := []struct {
		c         *Collector
		name      string
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
				tt.c.CollectNewMetrics()
			}
			expectedKeys := []string{
				"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
				"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
				"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
				"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
				"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
				"Sys", "TotalAlloc", "TotalMemory", "FreeMemory",
			}
			for _, key := range expectedKeys {
				value, ok := tt.c.Metrics.Load(key)
				require.True(t, ok, "Key %q should be present in metrics", key)
				require.NotZero(t, value, "Metric %q should not be zero", key)
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

			metricsMap := tt.c.GetAllMetrics()
			for _, key := range expectedKeys {
				value, ok := metricsMap[key]
				require.True(t, ok, "Key %q should be present in metricsMap", key)
				require.NotZero(t, value, "Metric %q should not be zero", key)
			}
			ctx, cancel := context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go tt.c.Collect(ctx, wg, 1*time.Millisecond, tt.c.CollectMetrics)
			time.Sleep(100 * time.Millisecond)
			cancel()
		})
	}
}
