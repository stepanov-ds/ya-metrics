package storage

import "github.com/stepanov-ds/ya-metrics/internal/utils"

type Storage interface {
	GetMetric(key string) (utils.Metric, bool)
	GetAllMetrics() map[string]utils.Metric
	SetMetricGauge(key string, value float64)
	SetMetricCounter(key string, value int64)
}
