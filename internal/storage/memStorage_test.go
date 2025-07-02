package storage

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

func TestMemStorage_GetMetric_NotFound(t *testing.T) {
	storage := NewMemStorage(&sync.Map{})
	metric, found := storage.GetMetric("unknown_key")
	require.False(t, found)
	require.Equal(t, utils.Metrics{}, metric)
}

func TestMemStorage_SetAndGet_Gauge(t *testing.T) {
	storage := NewMemStorage(&sync.Map{})
	key := "test_gauge"
	value := 3.14

	storage.SetMetric(context.Background(), key, value, false)

	metric, found := storage.GetMetric(key)
	require.True(t, found)
	assert.Equal(t, "gauge", metric.MType)
	assert.InDelta(t, 3.14, metric.Get(), 0.001)
}

func TestMemStorage_SetAndGet_Counter(t *testing.T) {
	storage := NewMemStorage(&sync.Map{})
	key := "test_counter"
	value := 42

	storage.SetMetric(context.Background(), key, value, true)

	metric, found := storage.GetMetric(key)
	require.True(t, found)
	assert.Equal(t, "counter", metric.MType)
	assert.Equal(t, int64(42), metric.Get())

	storage.SetMetric(context.Background(), key, 10, true)

	metric, found = storage.GetMetric(key)
	require.True(t, found)
	assert.Equal(t, int64(52), metric.Get())
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	storage := NewMemStorage(&sync.Map{})

	storage.SetMetric(context.Background(), "g1", 1.1, false)
	storage.SetMetric(context.Background(), "c1", 100, true)
	storage.SetMetric(context.Background(), "c2", 50, true)

	all := storage.GetAllMetrics()

	assert.Len(t, all, 3)

	assert.Contains(t, all, "g1")
	assert.Equal(t, "gauge", all["g1"].MType)
	assert.InDelta(t, 1.1, *all["g1"].Value, 0.001)

	assert.Contains(t, all, "c1")
	assert.Equal(t, "counter", all["c1"].MType)
	assert.Equal(t, int64(100), *all["c1"].Delta)

	assert.Contains(t, all, "c2")
	assert.Equal(t, "counter", all["c2"].MType)
	assert.Equal(t, int64(50), *all["c2"].Delta)
}

func TestMemStorage_OverwriteExistingGauge(t *testing.T) {
	storage := NewMemStorage(&sync.Map{})
	key := "temp"

	storage.SetMetric(context.Background(), key, 10.5, false)
	metric, _ := storage.GetMetric(key)
	assert.InDelta(t, 10.5, metric.Get(), 0.001)

	storage.SetMetric(context.Background(), key, 20.3, false)
	metric, _ = storage.GetMetric(key)
	assert.InDelta(t, 20.3, metric.Get(), 0.001)
}
