package metric

type Metric struct {
	Counter   int64
	Gauge     float64
	IsCounter bool
}

type MemStorage struct {
	Storage map[string]Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
        Storage: make(map[string]Metric),
    }
}