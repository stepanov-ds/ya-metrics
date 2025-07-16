package utils

import (
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func TestNewConstantIncreaseBackOff_InitialValues(t *testing.T) {
	b := NewConstantIncreaseBackOff(100*time.Millisecond, 200*time.Millisecond, 5)

	if b.initialInterval != 100*time.Millisecond {
		t.Errorf("Expected initial interval 100ms, got %v", b.initialInterval)
	}
	if b.currentInterval != 100*time.Millisecond {
		t.Errorf("Expected current interval 100ms, got %v", b.currentInterval)
	}
	if b.increase != 200*time.Millisecond {
		t.Errorf("Expected increase 200ms, got %v", b.increase)
	}
	if b.maxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", b.maxRetries)
	}
	if b.retry != 0 {
		t.Errorf("Expected retry count 0, got %d", b.retry)
	}
}

func TestNextBackOff_IncreasesCorrectly(t *testing.T) {
	b := NewConstantIncreaseBackOff(1*time.Second, 2*time.Second, 3)

	tests := []struct {
		expected time.Duration
	}{
		{1 * time.Second},
		{1*time.Second + 2*time.Second},
		{1*time.Second + 4*time.Second},
		{backoff.Stop},
	}

	for i, test := range tests {
		got := b.NextBackOff()
		if got != test.expected {
			t.Errorf("Test step %d: expected %v, got %v", i+1, test.expected, got)
		}
	}
}

func TestReset_RestoresInitialState(t *testing.T) {
	b := NewConstantIncreaseBackOff(1*time.Second, 2*time.Second, 3)

	b.NextBackOff()
	b.NextBackOff()

	b.Reset()

	if b.currentInterval != 1*time.Second {
		t.Errorf("Expected currentInterval to be reset to 1s, got %v", b.currentInterval)
	}
	if b.retry != 0 {
		t.Errorf("Expected retry count to be reset to 0, got %d", b.retry)
	}
}

func TestMaxRetries_Exceeded(t *testing.T) {
	b := NewConstantIncreaseBackOff(1*time.Second, 1*time.Second, 2)
	println(b.retry)
	if b.NextBackOff() != 1*time.Second {
		t.Errorf("First retry should return 1, got %v", b.currentInterval)
	}

	if b.NextBackOff() != 2*time.Second {
		t.Errorf("Second retry should return 2s, got %v", b.currentInterval)
	}

	if b.NextBackOff() != backoff.Stop {
		t.Errorf("Third retry should return backoff.Stop, got %v", b.currentInterval)
	}
}

func TestNewOneThreeFiveBackOff(t *testing.T) {
	b := NewOneThreeFiveBackOff()

	if b.initialInterval != time.Second {
		t.Errorf("Expected initial interval 1s, got %v", b.initialInterval)
	}
	if b.increase != 2*time.Second {
		t.Errorf("Expected increase 2s, got %v", b.increase)
	}
	if b.maxRetries != 3 {
		t.Errorf("Expected max retries 3, got %d", b.maxRetries)
	}

	if b.NextBackOff() != 1*time.Second {
		t.Errorf("Expected 1s for first retry, got %v", b.currentInterval)
	}
	if b.NextBackOff() != 3*time.Second {
		t.Errorf("Expected 3s for second retry, got %v", b.currentInterval)
	}
	if b.NextBackOff() != 5*time.Second {
		t.Errorf("Expected 5s for second retry, got %v", b.currentInterval)
	}
	if b.NextBackOff() != backoff.Stop {
		t.Errorf("Expected Stop after 3 retries, got %v", b.currentInterval)
	}
}
