package collector

import (
	"testing"

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
			c:         NewCollector(),
			pollCount: 1,
		},
		{
			name:      "Positive #2 collect twice",
			c:         NewCollector(),
			pollCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < int(tt.pollCount); i++ {
				tt.c.CollectMetrics()
			}
			assert.Equal(t, tt.pollCount, tt.c.Metrics["PollCount"].Counter)
		})
	}
}
