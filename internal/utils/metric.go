package utils

type Metric struct {
	Counter   int64   `json:"Counter,omitempty"`
	Gauge     float64 `json:"Gauge,omitempty"`
	IsCounter bool    `json:"-"`
}

func NewMetricCounter(counter int64) Metric {
	return Metric{
		Counter:   counter,
		IsCounter: true,
	}
}
func NewMetricGauge(gauge float64) Metric {
	return Metric{
		Gauge:     gauge,
		IsCounter: false,
	}
}
