package assistant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
)

func Test_clientProviderImpl_GetClient(t *testing.T) {
	tests := []struct {
		name               string
		config             Config
		expectedError      bool
		expectedClientType interface{}
	}{
		{
			name:               "openai",
			config:             Config{OpenAI: &openai.Config{Enabled: true}},
			expectedError:      false,
			expectedClientType: &openai.Client{},
		},
		{
			name:               "gemini",
			config:             Config{Gemini: &gemini.Config{Enabled: true}},
			expectedError:      false,
			expectedClientType: &gemini.Client{},
		},
		{
			name:          "multiple enabled - error",
			config:        Config{OpenAI: &openai.Config{Enabled: true}, Gemini: &gemini.Config{Enabled: true}},
			expectedError: true,
		},
		{
			name:          "none enabled - error",
			config:        Config{OpenAI: &openai.Config{Enabled: false}, Gemini: &gemini.Config{Enabled: false}},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := clientProviderImpl{}.GetClient(tt.config)
			if tt.expectedError {
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, client)
				assert.IsType(t, tt.expectedClientType, client)
			}
		})
	}
}
