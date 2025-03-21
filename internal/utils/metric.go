package utils

import (
	"fmt"
	"reflect"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetrics(name string, value interface{}, counter bool) Metrics {
	m := Metrics{
		ID: name,
	}
	m.Set(value, counter)
	return m
}

func (m *Metrics) Get() interface{} {
	if m.MType == "counter" {
		return *m.Delta
	} else if m.MType == "gauge" {
		return *m.Value
	} else {
		return nil
	}
}

func (m *Metrics) Set(value interface{}, counter bool) {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if counter {
		m.MType = "counter"
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if m.Delta == nil {
				m.Delta = new(int64)
				*m.Delta = v.Int()
			} else {
				*m.Delta = *m.Delta + v.Int()
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if m.Delta == nil {
				m.Delta = new(int64)
				*m.Delta = int64(v.Uint())
			} else {
				*m.Delta = *m.Delta + int64(v.Uint())
			}
		}
	} else {
		m.Value = new(float64)
		m.MType = "gauge"
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			*m.Value = float64(v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			*m.Value = float64(v.Uint())
		case reflect.Float32, reflect.Float64:
			*m.Value = v.Float()
		}
	}
}
func (m *Metrics) ConstructPath() (string, bool) {
	if m.MType == "counter" {
		return fmt.Sprintf("/update/counter/%s/%d", m.ID, *m.Delta), true
	} else if m.MType == "gauge" {
		return fmt.Sprintf("/update/gauge/%s/%f", m.ID, *m.Value), true
	} else {
		return "", false
	}
}
