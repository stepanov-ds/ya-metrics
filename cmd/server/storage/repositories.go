package storage

type Metric struct {
	Counter   int64
	Gauge     float64
	IsCounter bool
}

type Repositories interface {
	GetMetric(key string) (Metric)
	SetMetric(key string, m Metric)
}