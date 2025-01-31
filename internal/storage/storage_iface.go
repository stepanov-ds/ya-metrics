package storage

import "github.com/stepanov-ds/ya-metrics/internal/utils"


type Storage interface {
	GetMetric(key string) (utils.Metric, bool)
	SetMetric(key string, m utils.Metric)
	LockMutex()
	UnlockMutex()
	GetAllMetrics() map[string]utils.Metric
}
