package strutil

import (
	"testing"
)

func TestHardWrap(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		limit    int
		expected string
	}{
		{
			name:     "Short string",
			input:    "hello",
			limit:    10,
			expected: "hello",
		},
		{
			name:     "Exact limit",
			input:    "hello",
			limit:    5,
			expected: "hello",
		},
		{
			name:     "Wrap needed",
			input:    "helloworld",
			limit:    5,
			expected: "hello\nworld",
		},
		{
			name:     "Multiple wraps",
			input:    "helloworldagain",
			limit:    5,
			expected: "hello\nworld\nagain",
		},
		{
			name:     "Unicode wrap",
			input:    "こんにちは世界",
			limit:    2,
			expected: "こん\nにち\nは世\n界",
		},
		{
			name:     "Empty string",
			input:    "",
			limit:    5,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HardWrap(tt.input, tt.limit); got != tt.expected {
				t.Errorf("HardWrap() = %q, want %q", got, tt.expected)
			}
		})
	}
}
