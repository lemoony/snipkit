package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoDiscoveryConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    *Config
		expected Config
	}{
		{name: "nil config", input: nil, expected: Config{Enabled: true}},
		{name: "non nil config", input: &Config{Enabled: false}, expected: Config{Enabled: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := AutoDiscoveryConfig(tt.input)
			assert.Equal(t, config.Enabled, tt.expected.Enabled)
		})
	}
}
