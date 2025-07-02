// Package utils contains utility functions and shared types used across the application.
//
// This file defines the Metrics type and related methods for handling metric values.
package utils

import (
	"fmt"
	"reflect"
)

// Metrics represents a single metric with name, type, and value.
//
// Used for both gauge and counter types.
type Metrics struct {
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

// NewMetrics creates a new Metrics instance and sets its value based on type.
func NewMetrics(name string, value interface{}, counter bool) Metrics {
	m := Metrics{
		ID: name,
	}
	m.Set(value, counter)
	return m
}

// Get returns the current value of the metric.
//
// Returns:
// - int64 if it's a counter
// - float64 if it's a gauge
// - nil if type is unknown
func (m *Metrics) Get() interface{} {
	switch m.MType {
	case "counter":
		return *m.Delta
	case "gauge":
		return *m.Value
	default:
		return nil
	}
}

// Set updates the metric value according to its type.
//
// Supports various numeric types for both counter and gauge.
// For counters: increments existing value or sets a new one.
// For gauges: replaces current value.
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

// ConstructPath builds a URL path for updating this metric.
//
// Returns:
// - URL path string
// - true if successful, false if metric type is unknown
func (m *Metrics) ConstructPath() (string, bool) {
	switch m.MType {
	case "counter":
		return fmt.Sprintf("/update/counter/%s/%v", m.ID, *m.Delta), true
	case "gauge":
		return fmt.Sprintf("/update/gauge/%s/%v", m.ID, *m.Value), true
	default:
		return "", false
	}
}
