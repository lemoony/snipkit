package assistant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config_ValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError error
	}{
		{
			name:        "empty providers - valid",
			config:      Config{},
			expectError: nil,
		},
		{
			name:        "one enabled - valid",
			config:      Config{Providers: []ProviderConfig{{Type: ProviderTypeOpenAI, Enabled: true}}},
			expectError: nil,
		},
		{
			name:        "one disabled - valid",
			config:      Config{Providers: []ProviderConfig{{Type: ProviderTypeOpenAI, Enabled: false}}},
			expectError: nil,
		},
		{
			name: "multiple enabled - invalid",
			config: Config{Providers: []ProviderConfig{
				{Type: ProviderTypeOpenAI, Enabled: true},
				{Type: ProviderTypeGemini, Enabled: true},
			}},
			expectError: ErrorMultipleProvidersEnabled,
		},
		{
			name: "one enabled one disabled - valid",
			config: Config{Providers: []ProviderConfig{
				{Type: ProviderTypeOpenAI, Enabled: true},
				{Type: ProviderTypeGemini, Enabled: false},
			}},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.ValidateConfig()
			assert.Equal(t, tt.expectError, err)
		})
	}
}

func Test_Config_GetActiveProvider_MultipleEnabled(t *testing.T) {
	config := Config{Providers: []ProviderConfig{
		{Type: ProviderTypeOpenAI, Enabled: true},
		{Type: ProviderTypeGemini, Enabled: true},
	}}

	provider, err := config.GetActiveProvider()
	assert.Nil(t, provider)
	assert.Equal(t, ErrorMultipleProvidersEnabled, err)
}
