// Package utils contains utility functions and types used across the application.
//
// This file provides a custom backoff strategy for retrying failed operations.
package utils

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// ConstantIncreaseBackOff implements a custom backoff strategy with increasing intervals.
//
// It starts at initialInterval and increases by 'increase' after each retry,
// up to maxRetries attempts.
type ConstantIncreaseBackOff struct {
	initialInterval time.Duration
	currentInterval time.Duration
	increase        time.Duration
	maxRetries      int
	retry           int
}

// Reset sets the backoff strategy to its initial state.
func (b *ConstantIncreaseBackOff) Reset() {
	b.currentInterval = b.initialInterval
}

// NextBackOff returns the interval to wait before the next retry.
//
// Returns backoff.Stop if maximum retries reached.
func (b *ConstantIncreaseBackOff) NextBackOff() time.Duration {
	b.retry = b.retry + 1
	if b.retry < b.maxRetries {
		b.currentInterval = b.currentInterval + b.increase
		return b.currentInterval
	}
	return backoff.Stop
}

// NewConstantIncreaseBackOff creates and returns a new ConstantIncreaseBackOff instance.
//
// Parameters:
// - initial: initial backoff duration
// - inc: amount to increase backoff after each attempt
// - retries: maximum number of retries
func NewConstantIncreaseBackOff(initial time.Duration, inc time.Duration, retries int) *ConstantIncreaseBackOff {
	return &ConstantIncreaseBackOff{
		initialInterval: initial,
		increase:        inc,
		currentInterval: initial,
		maxRetries:      retries,
		retry:           0,
	}
}

// NewOneThreeFiveBackOff returns a backoff strategy with intervals [1s, 3s, 5s].
//
// Suitable for short retry sequences with increasing delays.
func NewOneThreeFiveBackOff() *ConstantIncreaseBackOff {
	return NewConstantIncreaseBackOff(time.Second, time.Second*2, 3)
}
