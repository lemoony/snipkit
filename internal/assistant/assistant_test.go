package assistant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	assistantMocks "github.com/lemoony/snipkit/mocks/assistant/client"
)

type clientProviderMock struct {
	mock.Mock
}

func (m *clientProviderMock) GetClient(config Config) (Client, error) {
	args := m.Called(config)
	return args.Get(0).(Client), args.Error(1)
}

func Test_AssistantImpl_Query(t *testing.T) {
	sys := testutil.NewTestSystem()
	cfg := Config{
		Providers: []ProviderConfig{
			{Type: ProviderTypeOpenAI, Enabled: true, Model: "gpt-4o", APIKeyEnv: "TEST_KEY"},
		},
	}

	clientMock := assistantMocks.NewClient(t)
	clientMock.On(mockutil.Query, mock.Anything).Return(`#!/bin/sh
#
# Simple Contents
# Filename: simple-Contents.sh
#
echo "foo"`, nil)

	providerMock := clientProviderMock{}
	providerMock.On("GetClient", mock.Anything).Return(clientMock, nil)

	assistant := NewBuilder(sys, cfg, cache.New(sys), withClientProvider(&providerMock))
	if ok, _ := assistant.Initialize(); !ok {
		assert.Fail(t, "assistant failed to initialize")
	}

	parsed := assistant.Query("test prompt")
	assert.Equal(t, `#!/bin/sh
#
# Simple Contents
#
echo "foo"`, parsed.Contents)
	assert.Equal(t, "simple-Contents.sh", parsed.Filename)
}

func Test_AssistantDescriptions(t *testing.T) {
	asst := assistantImpl{}

	descriptions := asst.AssistantDescriptions(Config{
		Providers: []ProviderConfig{
			{Type: ProviderTypeOpenAI, Enabled: false},
			{Type: ProviderTypeGemini, Enabled: true},
		},
	})

	// Should have descriptions for all supported providers
	assert.Len(t, descriptions, len(SupportedProviders))

	// Find OpenAI and Gemini in the descriptions
	var openaiDesc, geminiDesc *model.AssistantDescription
	for i := range descriptions {
		if descriptions[i].Key == model.AssistantKey(ProviderTypeOpenAI) {
			openaiDesc = &descriptions[i]
		}
		if descriptions[i].Key == model.AssistantKey(ProviderTypeGemini) {
			geminiDesc = &descriptions[i]
		}
	}

	assert.NotNil(t, openaiDesc)
	assert.False(t, openaiDesc.Enabled)
	assert.NotNil(t, geminiDesc)
	assert.True(t, geminiDesc.Enabled)
}

func Test_ValidateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		config        Config
		expectedValid bool
	}{
		{name: "not valid - empty providers", config: Config{}, expectedValid: false},
		{
			name:          "not valid - no provider enabled",
			config:        Config{Providers: []ProviderConfig{{Type: ProviderTypeOpenAI, Enabled: false}}},
			expectedValid: false,
		},
		{
			name:          "valid - one provider enabled",
			config:        Config{Providers: []ProviderConfig{{Type: ProviderTypeOpenAI, Enabled: true, Model: "gpt-4o", APIKeyEnv: "TEST_KEY"}}},
			expectedValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			clientMock := assistantMocks.NewClient(t)
			providerMock := clientProviderMock{}
			if tt.expectedValid {
				providerMock.On("GetClient", mock.Anything).Return(clientMock, nil)
			} else {
				providerMock.On("GetClient", mock.Anything).Return((*assistantMocks.Client)(nil), ErrorNoAssistantEnabled)
			}
			asst := assistantImpl{config: tt.config, provider: &providerMock}
			valid, msg := asst.Initialize()
			assert.Equal(t, tt.expectedValid, valid)
			if !tt.expectedValid {
				assert.Equal(t, msg, uimsg.AssistantNoneEnabled())
			}
		})
	}
}

func Test_AutoConfig(t *testing.T) {
	tests := []struct {
		name             string
		initialConfig    Config
		key              model.AssistantKey
		expectedEnabled  ProviderType
		expectedDisabled []ProviderType
	}{
		{
			name:            "no assistant configured yet",
			initialConfig:   Config{},
			key:             model.AssistantKey(ProviderTypeGemini),
			expectedEnabled: ProviderTypeGemini,
		},
		{
			name: "gemini enabled, auto configure openai",
			initialConfig: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeGemini, Enabled: true},
				},
			},
			key:              model.AssistantKey(ProviderTypeOpenAI),
			expectedEnabled:  ProviderTypeOpenAI,
			expectedDisabled: []ProviderType{ProviderTypeGemini},
		},
		{
			name: "openai enabled, auto configure gemini",
			initialConfig: Config{
				Providers: []ProviderConfig{
					{Type: ProviderTypeOpenAI, Enabled: true},
				},
			},
			key:              model.AssistantKey(ProviderTypeGemini),
			expectedEnabled:  ProviderTypeGemini,
			expectedDisabled: []ProviderType{ProviderTypeOpenAI},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := testutil.NewTestSystem()
			a := NewBuilder(sys, tt.initialConfig, cache.New(sys))
			updateConfig := a.AutoConfig(tt.key)

			// Check that the expected provider is enabled
			var foundEnabled bool
			for _, p := range updateConfig.Providers {
				if p.Type == tt.expectedEnabled {
					assert.True(t, p.Enabled, "expected %s to be enabled", tt.expectedEnabled)
					foundEnabled = true
				}
			}
			assert.True(t, foundEnabled, "expected provider %s not found", tt.expectedEnabled)

			// Check that expected providers are disabled
			for _, disabledType := range tt.expectedDisabled {
				for _, p := range updateConfig.Providers {
					if p.Type == disabledType {
						assert.False(t, p.Enabled, "expected %s to be disabled", disabledType)
					}
				}
			}
		})
	}
}
