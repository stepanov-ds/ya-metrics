package utils

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

type ConstantIncreaseBackOff struct {
	initialInterval time.Duration
	currentInterval time.Duration
	increase        time.Duration
	maxRetries      int
	retry           int
}

func (b *ConstantIncreaseBackOff) Reset() {
	b.currentInterval = b.initialInterval
}
func (b *ConstantIncreaseBackOff) NextBackOff() time.Duration {
	b.retry = b.retry + 1
	if b.retry < b.maxRetries {
		b.currentInterval = b.currentInterval + b.increase
		return b.currentInterval
	}
	return backoff.Stop
}

func NewConstantIncreaseBackOff(initial time.Duration, inc time.Duration, retries int) *ConstantIncreaseBackOff {
	return &ConstantIncreaseBackOff{
		initialInterval: initial,
		increase:        inc,
		currentInterval: initial,
		maxRetries:      retries,
		retry:           0,
	}
}
