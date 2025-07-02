package utils

import (
	"testing"
)

func TestCalculateHashWithKey(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		key      string
		expected string
	}{
		{
			name:     "Standard input",
			data:     "testdata",
			key:      "secretkey",
			expected: "e06638745371199b52204e6bbdaeef89cd04d47e38c1e7d63eddd6fbb7329857",
		},
		{
			name:     "Empty data",
			data:     "",
			key:      "secretkey",
			expected: "0ec2fbe02ea7c3eb6dd73c12eb2cffc9061280dfc8365cdcfa5241c6e3d9c9a7",
		},
		{
			name:     "Empty key",
			data:     "testdata",
			key:      "",
			expected: "09e36bb2ed1154cdd3f5b5d17dedd9c727acb38a62924c36ee4a214ccf463b97",
		},
		{
			name:     "Numbers and symbols",
			data:     "1234567890!@#$%^&*()",
			key:      "mysecretpassword",
			expected: "7bfa8e5c4ff52f0f02982b470e44a60cdcac7eb1aa64066ade54cd3a21465f17",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CalculateHashWithKey([]byte(test.data), test.key)
			if result != test.expected {
				t.Errorf("Expected hash %s, got %s", test.expected, result)
			}
		})
	}
}
