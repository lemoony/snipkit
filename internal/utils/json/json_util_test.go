package jsonutil

import "testing"

func TestCompactJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "simple json with whitespace",
			input:    []byte(`{  "test":"foo"  }`),
			expected: `{"test":"foo"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := CompactJSON(tt.input); actual != tt.expected {
				t.Errorf("CompactJSON() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
