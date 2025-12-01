package strutil

import (
	"testing"
)

func TestTruncateMiddle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "Short string",
			input:    "hello",
			n:        3,
			expected: "hello",
		},
		{
			name:     "Exact length string",
			input:    "abcdef",
			n:        3,
			expected: "abcdef",
		},
		{
			name:     "Long string",
			input:    "hello world this is a long string",
			n:        5,
			expected: "hello...tring",
		},
		{
			name:     "Empty string",
			input:    "",
			n:        5,
			expected: "",
		},
		{
			name:     "Unicode string",
			input:    "こんにちは世界",
			n:        2,
			expected: "こん...世界",
		},
		{
			name:     "Unicode short string",
			input:    "こんにちは",
			n:        3,
			expected: "こんにちは",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateMiddle(tt.input, tt.n); got != tt.expected {
				t.Errorf("TruncateMiddle() = %v, want %v", got, tt.expected)
			}
		})
	}
}
