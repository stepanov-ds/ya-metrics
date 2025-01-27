package storage

type MemStorage struct {
	storage map[string]Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
        storage: make(map[string]Metric),
    }
}

func (s *MemStorage) GetMetric(key string) Metric{
	return s.storage[key]
}

func (s *MemStorage) SetMetric(key string, m Metric) {
	s.storage[key] = m
}