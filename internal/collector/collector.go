package collector

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

type Collector struct {
	Metrics *sync.Map
}

func NewCollector(m *sync.Map) *Collector {
	return &Collector{
		Metrics: m,
	}
}

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
		c.Metrics.Store(k, utils.NewMetricGauge(v))
	}

	if value, exist := c.Metrics.Load("PollCount"); exist {
		if v, ok := value.(*utils.MetricCounter); ok {
			v.Set(1)
		} else {
			c.Metrics.Store("PollCount", utils.NewMetricCounter(1))
		}
	} else {
		c.Metrics.Store("PollCount", utils.NewMetricCounter(1))
	}

	c.Metrics.Store("RandomValue", utils.NewMetricGauge(rand.Float64()))

}
func (c *Collector) GetAllMetrics() map[string]utils.Metric {
	result := make(map[string]utils.Metric)

	c.Metrics.Range(func(key, value interface{}) bool {
		k, ok1 := key.(string)
		v, ok2 := value.(utils.Metric)

		if ok1 && ok2 {
			result[k] = v
		}
		return true
	})

	return result
}

func (c *Collector) collect(interval time.Duration) {
	for {
		c.CollectMetrics()
		time.Sleep(interval)
	}
}

func (c *Collector) Collect(interval time.Duration) {
	go c.collect(interval)
}
