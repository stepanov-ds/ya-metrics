package storage

import (
	"context"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

type Storage interface {
	GetMetric(key string) (utils.Metrics, bool)
	GetAllMetrics() map[string]utils.Metrics
	SetMetric(ctx context.Context, key string, value interface{}, counter bool)
}
