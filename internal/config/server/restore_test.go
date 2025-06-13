package server

import (
	"context"
	"math/rand/v2"
	"os"
	"sync"
	"testing"

	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestRestore(t *testing.T) {
	*FileStorePath = "filestore_test.out"
	expectedMetrics := map[string]bool{
		"HeapAlloc":  false,
		"StackInuse": false,
		"PollCount":  true,
	}
	defer os.Remove(*FileStorePath)
	st := storage.NewMemStorage(&sync.Map{})

	for k, v := range expectedMetrics {
		st.SetMetric(context.Background(), k, rand.IntN(100), v)
	}

	storeInFile(st)
	restored := RestoreStorage()

	for k, _ := range expectedMetrics {
		metricOrigin, ok := st.GetMetric(k)
		assert.True(t, ok, "metric not found in origin storage")
		metricRestored, ok := restored.GetMetric(k)
		assert.True(t, ok, "metric not found in restored storage")
		assert.Equal(t, metricOrigin, metricRestored)
	}
}
