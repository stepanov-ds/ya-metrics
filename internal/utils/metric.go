package utils

import (
	"fmt"
	"reflect"
)

type Metric interface {
	Get() interface{}
	Set(interface{})
	ConstructPath(string) string
	ConstructJsonObj(string) Metrics
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricCounter struct {
	Counter int64 `json:"delta,omitempty"`
}

type MetricGauge struct {
	Gauge float64 `json:"value,omitempty"`
}

func NewMetricCounter(counter int64) *MetricCounter {
	return &MetricCounter{
		Counter: counter,
	}
}
func NewMetricGauge(gauge float64) *MetricGauge {
	return &MetricGauge{
		Gauge: gauge,
	}
}

func (m *MetricCounter) Get() interface{} {
	return m.Counter
}

func (m *MetricCounter) Set(value interface{}) {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		m.Counter = m.Counter + v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		m.Counter = m.Counter + int64(v.Uint())
	}
}
func (m *MetricCounter) ConstructPath(name string) string {
	return fmt.Sprintf("/update/counter/%s/%d", name, m.Counter)
}
func (m *MetricGauge) Get() interface{} {
	return m.Gauge
}
func (m *MetricGauge) Set(value interface{}) {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		m.Gauge = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		m.Gauge = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		m.Gauge = v.Float()
	}
}
func (m *MetricGauge) ConstructPath(name string) string {
	return fmt.Sprintf("/update/gauge/%s/%f", name, m.Gauge)
}

func (m *MetricGauge) ConstructJsonObj(name string) Metrics {
	return Metrics{
		ID:    name,
		MType: "gauge",
		Value: &m.Gauge,
	}
}

func (m *MetricCounter) ConstructJsonObj(name string) Metrics {
	return Metrics{
		ID:    name,
		MType: "counter",
		Delta: &m.Counter,
	}
}
