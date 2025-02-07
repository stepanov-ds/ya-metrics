package utils

import (
	"fmt"
	"reflect"
)

type Metric interface {
	Get() interface{}
	Set(interface{})
	ConstructPath(string) string
}
type MetricCounter struct {
	Counter int64 `json:"Counter,omitempty"`
}

type MetricGauge struct {
	Gauge float64 `json:"Gauge,omitempty"`
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
