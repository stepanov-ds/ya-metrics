package utils

import (
	"testing"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics("test_gauge", 3.14, false)
	if m.ID != "test_gauge" {
		t.Errorf("Expected ID 'test_gauge', got '%s'", m.ID)
	}
	if m.MType != "gauge" {
		t.Errorf("Expected type 'gauge', got '%s'", m.MType)
	}
	if m.Get() != 3.14 {
		t.Errorf("Expected value 3.14, got %v", m.Get())
	}
}

func TestSetAndGet_Gauge(t *testing.T) {
	var m Metrics
	m.Set(42.5, false)
	if m.MType != "gauge" {
		t.Errorf("Expected type 'gauge', got '%s'", m.MType)
	}
	if val, ok := m.Get().(float64); ok {
		if val != 42.5 {
			t.Errorf("Expected 42.5, got %f", val)
		}
	} else {
		t.Error("Expected float64 for gauge")
	}
}

func TestSetAndGet_Counter(t *testing.T) {
	var m Metrics
	m.Set(10, true)
	m.Set(20, true) // increment

	if m.MType != "counter" {
		t.Errorf("Expected type 'counter', got '%s'", m.MType)
	}
	if val, ok := m.Get().(int64); ok {
		if val != 30 {
			t.Errorf("Expected 30, got %d", val)
		}
	} else {
		t.Error("Expected int64 for counter")
	}
}

func TestConstructPath_Gauge(t *testing.T) {
	m := Metrics{
		ID:    "some_metric",
		MType: "gauge",
		Value: new(float64),
	}
	*m.Value = 123.45

	path, ok := m.ConstructPath()
	if !ok {
		t.Error("Expected path to be constructed for gauge")
	}
	expected := "/update/gauge/some_metric/123.45"
	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

func TestConstructPath_Counter(t *testing.T) {
	m := Metrics{
		ID:    "clicks",
		MType: "counter",
		Delta: new(int64),
	}
	*m.Delta = 42

	path, ok := m.ConstructPath()
	if !ok {
		t.Error("Expected path to be constructed for counter")
	}
	expected := "/update/counter/clicks/42"
	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

func TestSet_WithDifferentTypes(t *testing.T) {
	tests := []struct {
		input    any
		expected any
		name     string
		counter  bool
	}{
		{42, float64(42), "int", false},
		{uint(100), float64(100), "uint", false},
		{float64(3.14), float64(3.14), "float32", false},
		{int64(99), int64(99), "int64", true},
		{uint64(200), int64(200), "uint64", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var m Metrics
			m.Set(test.input, test.counter)
			got := m.Get()

			if test.counter {
				if val, ok := got.(int64); !ok || val != test.expected.(int64) {
					t.Errorf("Expected %v (%T), got %v (%T)", test.expected, test.expected, got, got)
				}
			} else {
				if val, ok := got.(float64); !ok || val != test.expected.(float64) {
					t.Errorf("Expected %v (%T), got %v (%T)", test.expected, test.expected, got, got)
				}
			}
		})
	}
}
