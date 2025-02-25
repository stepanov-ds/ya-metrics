package storage

import "github.com/stepanov-ds/ya-metrics/internal/utils"

type Storage interface {
	GetMetric(key string) (utils.Metrics, bool)
	GetAllMetrics() map[string]utils.Metrics
	SetMetric(key string, value interface{}, counter bool)
}
