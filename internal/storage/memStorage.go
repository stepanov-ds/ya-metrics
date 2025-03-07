package storage

import (
	"sync"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

type MemStorage struct {
	storage *sync.Map
}

func NewMemStorage(m *sync.Map) *MemStorage {
	return &MemStorage{
		storage: m,
	}
}

func (s *MemStorage) GetMetric(key string) (utils.Metrics, bool) {
	value, found := s.storage.Load(key)
	if !found {
		return utils.Metrics{}, false
	}
	metric, ok := value.(utils.Metrics)
	if !ok {
		return utils.Metrics{}, false
	}
	return metric, true
}

func (s *MemStorage) SetMetric(key string, value interface{}, counter bool) {
	oldMetricValue, found := s.storage.Load(key)
	if found {
		switch v := oldMetricValue.(type) {
		case utils.Metrics:
			v.Set(value, counter)
			s.storage.Store(key, v)
		}
	} else {
		s.storage.Store(key, utils.NewMetrics(key, value, counter))
	}
}

func (s *MemStorage) GetAllMetrics() map[string]utils.Metrics {
	result := make(map[string]utils.Metrics)
	s.storage.Range(func(key, value interface{}) bool {
		if metric, ok := value.(utils.Metrics); ok {
			result[key.(string)] = metric
		}
		return true
	})
	return result
}
