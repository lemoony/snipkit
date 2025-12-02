package assistant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/assistant/langchain"
)

func Test_clientProviderImpl_GetClient(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		expectedError bool
	}{
		{
			name: "openai provider enabled",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: true, Model: "gpt-4o", APIKeyEnv: "TEST_KEY"},
				},
			},
			// Will error because API key env var is not set in test environment
			expectedError: true,
		},
		{
			name: "gemini provider enabled",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeGemini, Enabled: true, Model: "gemini-1.5-flash", APIKeyEnv: "TEST_KEY"},
				},
			},
			// Will error because API key env var is not set in test environment
			expectedError: true,
		},
		{
			name: "none enabled - error",
			config: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: false},
					{Type: ProviderTypeGemini, Enabled: false},
				},
			},
			expectedError: true,
		},
		{
			name:          "empty providers - error",
			config:        Config{},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := clientProviderImpl{}.GetClient(tt.config)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.IsType(t, &langchain.Client{}, client)
			}
		})
	}
}
