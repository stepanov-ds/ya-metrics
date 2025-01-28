package storage

import "github.com/stepanov-ds/ya-metrics/pkg/utils"

type Repositories interface {
	GetMetric(key string) utils.Metric
	SetMetric(key string, m utils.Metric)
}
