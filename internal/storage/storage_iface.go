// Package storage provides interfaces and implementations for storing and retrieving metrics.
//
// This file defines the core Storage interface used throughout the application.
package storage

import (
	"context"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

// Storage is the primary interface for metric persistence layer.
//
// Implementations must provide thread-safe access.
type Storage interface {
	// GetMetric retrieves a metric by key.
	// Returns the metric and true if found, empty metric and false otherwise.
	GetMetric(key string) (utils.Metrics, bool)
	// GetAllMetrics returns all stored metrics as a map of name to value.
	// Should return only valid metrics.
	GetAllMetrics() map[string]utils.Metrics
	// SetMetric stores or updates a metric with given type (counter/gauge).
	SetMetric(ctx context.Context, key string, value interface{}, counter bool)
}
