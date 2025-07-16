package utils

import (
	"os"
	"testing"
)

func TestGetFlagValue(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected string
		args     []string
	}{
		{
			name:     "separate flag",
			args:     []string{"cmd", "-k", "value"},
			flagName: "k",
			expected: "value",
		},
		{
			name:     "inline flag",
			args:     []string{"cmd", "-k=value"},
			flagName: "k",
			expected: "value",
		},
		{
			name:     "flag not found",
			args:     []string{"cmd", "-other", "value"},
			flagName: "k",
			expected: "",
		},
		{
			name:     "flag without value",
			args:     []string{"cmd", "-k"},
			flagName: "k",
			expected: "",
		},
		{
			name:     "partly named flag",
			args:     []string{"cmd", "-key=value"},
			flagName: "k",
			expected: "",
		},
		{
			name:     "no args",
			args:     []string{},
			flagName: "k",
			expected: "",
		},
		{
			name:     "long name flag",
			args:     []string{"cmd", "-config=value"},
			flagName: "config",
			expected: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			result := GetFlagValue(tt.flagName)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
