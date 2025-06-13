// Package storage provides interfaces and implementations for storing and retrieving metrics.
//
// This file contains MemStorage — an in-memory implementation based on sync.Map.
package storage

import (
	"context"
	"sync"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

// MemStorage is an in-memory implementation of the Storage interface.
//
// Uses sync.Map for thread-safe operations.
type MemStorage struct {
	storage *sync.Map
}

// NewMemStorage creates a new in-memory storage backed by the provided sync.Map.
func NewMemStorage(m *sync.Map) *MemStorage {
	return &MemStorage{
		storage: m,
	}
}

// GetMetric retrieves a metric by its key from the in-memory storage.
//
// Returns the metric and true if found, empty metric and false otherwise.
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

// SetMetric stores or updates a metric in memory.
//
// If the metric exists, it updates the value using Set method.
// If not, creates a new metric with given value and type.
func (s *MemStorage) SetMetric(ctx context.Context, key string, value interface{}, counter bool) {
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

// GetAllMetrics returns all stored metrics as a map[string]utils.Metrics.
//
// Returns only valid Metrics values, skipping any invalid entries.
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
