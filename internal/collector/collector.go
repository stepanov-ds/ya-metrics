// Package collector implements logic for collecting application metrics.
//
// Collector provides methods to gather runtime statistics including:
// - Memory usage (heap, GC)
// - CPU utilization
// - Virtual memory stats
// - Random values and counters
package collector

import (
	"context"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

// Collector is a structure responsible for gathering and storing metrics.
type Collector struct {
	Metrics *sync.Map // Metrics storage in the form map[metricName]metricValue
}

// NewCollector creates a new instance of Collector.
//
//	m - pointer to a sync.Map where metrics will be stored
func NewCollector(m *sync.Map) *Collector {
	return &Collector{
		Metrics: m,
	}
}

// CollectMetrics gathers memory and garbage collection metrics from the runtime.
// It updates the following metrics:
// - Alloc, BuckHashSys, Frees, GCCPUFraction, GCSys, HeapAlloc, HeapIdle etc.
// Also increases the "PollCount" counter and sets a new "RandomValue".
func (c *Collector) CollectMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics := map[string]float64{
		"Alloc":         float64(m.Alloc),
		"BuckHashSys":   float64(m.BuckHashSys),
		"Frees":         float64(m.Frees),
		"GCCPUFraction": m.GCCPUFraction,
		"GCSys":         float64(m.GCSys),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"HeapSys":       float64(m.HeapSys),
		"LastGC":        float64(m.LastGC),
		"Lookups":       float64(m.Lookups),
		"MCacheInuse":   float64(m.MCacheInuse),
		"MCacheSys":     float64(m.MCacheSys),
		"MSpanInuse":    float64(m.MSpanInuse),
		"MSpanSys":      float64(m.MSpanSys),
		"Mallocs":       float64(m.Mallocs),
		"NextGC":        float64(m.NextGC),
		"NumForcedGC":   float64(m.NumForcedGC),
		"NumGC":         float64(m.NumGC),
		"OtherSys":      float64(m.OtherSys),
		"PauseTotalNs":  float64(m.PauseTotalNs),
		"StackInuse":    float64(m.StackInuse),
		"StackSys":      float64(m.StackSys),
		"Sys":           float64(m.Sys),
		"TotalAlloc":    float64(m.TotalAlloc),
	}

	for k, v := range metrics {
		c.Metrics.Store(k, utils.NewMetrics(k, v, false))
	}

	if value, exist := c.Metrics.Load("PollCount"); exist {
		if v, ok := value.(utils.Metrics); ok {
			v.Set(1, true)
		}
	} else {
		c.Metrics.Store("PollCount", utils.NewMetrics("PollCount", 1, true))
	}

	c.Metrics.Store("RandomValue", utils.NewMetrics("RandomValue", rand.Float64(), false))

}

// CollectNewMetrics gathers extended system metrics like virtual memory and CPU utilization.
// It stores:
// - TotalMemory: total system memory
// - FreeMemory: free memory
// - CPUutilizationX: per-core CPU usage percentage
func (c *Collector) CollectNewMetrics() {
	v, _ := mem.VirtualMemory()
	if v != nil {
		c.Metrics.Store("TotalMemory", utils.NewMetrics("TotalMemory", float64(v.Total), false))
		c.Metrics.Store("FreeMemory", utils.NewMetrics("FreeMemory", float64(v.Free), false))
	}

	cpuUtil, _ := cpu.Percent(0, true)
	for i, v := range cpuUtil {
		c.Metrics.Store("CPUutilization"+strconv.FormatInt(int64(i), 10), utils.NewMetrics("CPUutilization"+strconv.FormatInt(int64(i), 10), v, false))
	}

}

// GetAllMetrics returns all stored metrics as a map.
// Returns:
// - map[string]utils.Metrics: a copy of all current metrics
func (c *Collector) GetAllMetrics() map[string]utils.Metrics {
	result := make(map[string]utils.Metrics)

	c.Metrics.Range(func(key, value interface{}) bool {
		k, ok1 := key.(string)
		v, ok2 := value.(utils.Metrics)

		if ok1 && ok2 {
			result[k] = v
		}
		return true
	})

	return result
}

// collect runs the given function f periodically at specified intervals.
// This is a helper method used by Collect().
func (c *Collector) collect(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, f func()) {
	defer wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f()
			time.Sleep(interval)
		}
	}
}

// Collect starts the metric collection loop in a separate goroutine.
// Calls the provided function f every interval.
func (c *Collector) Collect(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, f func()) {
	go c.collect(ctx, wg, interval, f)
}
