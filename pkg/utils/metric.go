package utils

type Metric struct {
	Counter   int64
	Gauge     float64
	IsCounter bool
}

func NewMetricCounter(counter int64) Metric {
	return Metric{
		Counter: counter,
		IsCounter: true,
	}
}
func NewMetricGauge(gauge float64) Metric {
	return Metric{
		Gauge: gauge,
		IsCounter: false,
	}
}
