// Package logger provides a global zap.Logger instance for structured logging.
//
// It allows initialization with custom log level and is used across the application
// to log events, errors, and operational information.
package logger

import "go.uber.org/zap"

// Log is a global logger instance initialized during application startup.
var Log *zap.Logger

// Initialize configures and builds a new zap.Logger with the specified log level.
//
// Uses zap.NewProductionConfig() as base configuration.
// Returns error if level parsing or logger creation fails.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}
