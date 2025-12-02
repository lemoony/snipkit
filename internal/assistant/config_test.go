package assistant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config_GetActiveProvider_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError error
		expectType  ProviderType
	}{
		{
			name:        "no providers - error",
			config:      Config{Providers: []ProviderConfig{}},
			expectError: ErrorNoAssistantEnabled,
		},
		{
			name: "all disabled - error",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: false},
					{Type: ProviderTypeGemini, Enabled: false},
				},
			},
			expectError: ErrorNoAssistantEnabled,
		},
		{
			name: "first enabled - success",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: true},
					{Type: ProviderTypeGemini, Enabled: false},
				},
			},
			expectType: ProviderTypeOpenAI,
		},
		{
			name: "second enabled - success",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: false},
					{Type: ProviderTypeGemini, Enabled: true},
				},
			},
			expectType: ProviderTypeGemini,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := tt.config.GetActiveProvider()

			if tt.expectError != nil {
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				assert.Equal(t, tt.expectType, provider.Type)
			}
		})
	}
}

func Test_Config_HasEnabledProvider(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name: "no enabled providers",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: false},
					{Type: ProviderTypeGemini, Enabled: false},
				},
			},
			expected: false,
		},
		{
			name: "has enabled provider",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: true},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.HasEnabledProvider()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_Config_ClientKey(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectedKey string
		expectError bool
	}{
		{
			name: "success - returns provider type as key",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: true},
				},
			},
			expectedKey: string(ProviderTypeOpenAI),
		},
		{
			name: "no enabled provider - error",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: false},
				},
			},
			expectError: true,
		},
		{
			name: "multiple enabled - error",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: true},
					{Type: ProviderTypeGemini, Enabled: true},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := tt.config.ClientKey()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedKey, string(key))
			}
		})
	}
}

func Test_DefaultProviderConfig_AllProviders(t *testing.T) {
	tests := []struct {
		providerType      ProviderType
		expectedModel     string
		expectedAPIKeyEnv string
		expectedServerURL string
	}{
		{
			providerType:      ProviderTypeOpenAI,
			expectedModel:     "gpt-4.1",
			expectedAPIKeyEnv: "SNIPKIT_OPENAI_API_KEY",
		},
		{
			providerType:      ProviderTypeAnthropic,
			expectedModel:     "claude-sonnet-4.5",
			expectedAPIKeyEnv: "SNIPKIT_ANTHROPIC_API_KEY",
		},
		{
			providerType:      ProviderTypeGemini,
			expectedModel:     "gemini-1.5-flash",
			expectedAPIKeyEnv: "SNIPKIT_GEMINI_API_KEY",
		},
		{
			providerType:      ProviderTypeOllama,
			expectedModel:     "llama3",
			expectedAPIKeyEnv: "",
			expectedServerURL: "http://localhost:11434",
		},
		{
			providerType:      ProviderTypeOpenAICompatible,
			expectedModel:     "",
			expectedAPIKeyEnv: "",
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.providerType), func(t *testing.T) {
			cfg := DefaultProviderConfig(tt.providerType)

			assert.Equal(t, tt.providerType, cfg.Type)
			assert.True(t, cfg.Enabled)
			assert.Equal(t, tt.expectedModel, cfg.Model)
			assert.Equal(t, tt.expectedAPIKeyEnv, cfg.APIKeyEnv)

			if tt.expectedServerURL != "" {
				assert.Equal(t, tt.expectedServerURL, cfg.ServerURL)
			}
		})
	}
}
