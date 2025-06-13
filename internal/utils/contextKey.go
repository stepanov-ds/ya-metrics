// Package utils contains utility functions and shared types used across the application.
//
// This file defines context keys for use in middleware and handlers.
package utils

// ContextKey is a custom type for keys used in context.WithValue
// to avoid collisions with other context values.
type ContextKey string

const (
	// Transaction is a context key for storing an active database transaction.
	Transaction ContextKey = "transaction"
)
