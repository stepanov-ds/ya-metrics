package storage

import (
	"sync"

	"github.com/stepanov-ds/ya-metrics/pkg/utils"
)

type MemStorage struct {
	storage map[string]utils.Metric
	mu      sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]utils.Metric),
	}
}

func (s *MemStorage) GetMetric(key string) (utils.Metric, bool) {
	m, found := s.storage[key]
	if found {
		return m, true
	} else {
		return utils.Metric{}, false
	}
}

func (s *MemStorage) SetMetric(key string, m utils.Metric) {
	if m.IsCounter{
		oldMetricValue, found := s.GetMetric(key)
		if found {
			if oldMetricValue.IsCounter {
				s.storage[key] = utils.Metric{
					Counter:   m.Counter + oldMetricValue.Counter,
					IsCounter: true,
				}
			} else {
				s.storage[key] = m
			}
		} else {
			s.storage[key] = m
		}
	} else {
		s.storage[key] = m
	}
	
	
}

func (s *MemStorage) LockMutex() {
	s.mu.Lock()
}
func (s *MemStorage) UnlockMutex() {
	s.mu.Unlock()
}
func (s *MemStorage) GetAllMetrics() map[string]utils.Metric {
	return s.storage
}
