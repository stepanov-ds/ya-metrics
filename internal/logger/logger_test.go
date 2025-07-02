package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize_WithValidLevels(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantPanic bool
	}{
		{"debug", "debug", false},
		{"info", "info", false},
		{"warn", "warn", false},
		{"error", "error", false},
		{"fatal", "fatal", false},
		{"invalid", "invalid_level", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			err := Initialize(tt.level)

			if tt.wantPanic {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unrecognized level")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, Log)
			}
		})
	}
}
