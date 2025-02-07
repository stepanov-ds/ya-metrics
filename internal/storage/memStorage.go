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

func (s *MemStorage) GetMetric(key string) (utils.Metric, bool) {
	value, found := s.storage.Load(key)
	if !found {
		return nil, false
	}
	metric, ok := value.(utils.Metric)
	if !ok {
		return nil, false
	}
	return metric, true
}

func (s *MemStorage) SetMetricGauge(key string, value float64) {
	s.storage.Store(key, utils.NewMetricGauge(value))
}

func (s *MemStorage) SetMetricCounter(key string, value int64) {
	oldMetricValue, found := s.storage.Load(key)
	if found {
		switch v := oldMetricValue.(type) {
		case *utils.MetricCounter:
			v.Set(value)
			s.storage.Store(key, v) // Обновляем значение в sync.Map
		case *utils.MetricGauge:
			s.storage.Store(key, utils.NewMetricCounter(value))
		}
	} else {
		s.storage.Store(key, utils.NewMetricCounter(value))
	}
}

func (s *MemStorage) GetAllMetrics() map[string]utils.Metric {
	result := make(map[string]utils.Metric)
	s.storage.Range(func(key, value interface{}) bool {
		if metric, ok := value.(utils.Metric); ok {
			result[key.(string)] = metric
		}
		return true
	})
	return result
}
