package storage

import "github.com/stepanov-ds/ya-metrics/pkg/utils"

type MemStorage struct {
	storage map[string]utils.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]utils.Metric),
	}
}

func (s *MemStorage) GetMetric(key string) utils.Metric {
	return s.storage[key]
}

func (s *MemStorage) SetMetric(key string, m utils.Metric) {
	s.storage[key] = m
}