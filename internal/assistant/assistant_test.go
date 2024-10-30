package assistant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/testutil"
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
		OpenAI: &openai.Config{Enabled: true},
	}

	clientMock := assistantMocks.NewClient(t)
	clientMock.On("Query", mock.Anything).Return(`#!/bin/sh
#
# Simple script
# Filename: simple-script.sh
#
echo "foo"`, nil)

	providerMock := clientProviderMock{}
	providerMock.On("GetClient", mock.Anything).Return(clientMock, nil)

	assistant := NewBuilder(sys, cfg, cache.New(sys), withClientProvider(&providerMock))

	// Assuming that the openai.NewClient and gemini.NewClient are mocked to return test clients.
	result1, result2 := assistant.Query("test prompt")
	assert.Equal(t, `#!/bin/sh
#
# Simple script
#
echo "foo"`, result1)
	assert.Equal(t, "simple-script.sh", result2)
}

func Test_AssistantDescriptions(t *testing.T) {
	asst := assistantImpl{}

	descriptions := asst.AssistantDescriptions(Config{
		OpenAI: &openai.Config{Enabled: false},
		Gemini: &gemini.Config{Enabled: true},
	})

	assert.Len(t, descriptions, 2)
	assert.Equal(t, openai.Key, descriptions[0].Key)
	assert.False(t, descriptions[0].Enabled)
	assert.Equal(t, gemini.Key, descriptions[1].Key)
	assert.True(t, descriptions[1].Enabled)
}

func Test_ValidateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		config        Config
		expectedValid bool
		panics        bool
	}{
		{"not valid", Config{}, false, false},
		{"more than one expected", Config{OpenAI: &openai.Config{Enabled: true}, Gemini: &gemini.Config{Enabled: true}}, false, true},
		{"valid", Config{OpenAI: &openai.Config{Enabled: true}}, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			asst := assistantImpl{config: tt.config, provider: clientProviderImpl{}}

			if tt.panics {
				assert.Panics(t, func() {
					_, _ = asst.ValidateConfig()
				})
			} else {
				valid, msg := asst.ValidateConfig()
				assert.Equal(t, tt.expectedValid, valid)
				if !tt.expectedValid {
					assert.Equal(t, msg, uimsg.AssistantNoneEnabled())
				}
			}
		})
	}
}

func Test_AutoConfig(t *testing.T) {
	tests := []struct {
		name          string
		initialConfig Config
		key           model.AssistantKey
		expected      Config
	}{
		{
			name:          "no assistant configured yet",
			initialConfig: Config{},
			key:           gemini.Key,
			expected: Config{
				Gemini: &gemini.Config{Enabled: true},
			},
		},
		{
			name: "gemini enabled, auto configure openai",
			initialConfig: Config{
				Gemini: &gemini.Config{Enabled: true},
			},
			key: openai.Key,
			expected: Config{
				OpenAI: &openai.Config{Enabled: true},
				Gemini: &gemini.Config{Enabled: false},
			},
		},
		{
			name: "openai enabled, auto configure gemini",
			initialConfig: Config{
				OpenAI: &openai.Config{Enabled: true},
			},
			key: gemini.Key,
			expected: Config{
				OpenAI: &openai.Config{Enabled: false},
				Gemini: &gemini.Config{Enabled: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := testutil.NewTestSystem()
			a := NewBuilder(sys, tt.initialConfig, cache.New(sys))
			updateConfig := a.AutoConfig(tt.key)

			assert.Equal(t, tt.expected.OpenAI != nil, updateConfig.OpenAI != nil)
			if tt.expected.OpenAI != nil {
				assert.Equal(t, tt.expected.OpenAI.Enabled, updateConfig.OpenAI.Enabled)
			}

			assert.Equal(t, tt.expected.Gemini != nil, updateConfig.Gemini != nil)
			if tt.expected.Gemini != nil {
				assert.Equal(t, tt.expected.Gemini.Enabled, updateConfig.Gemini.Enabled)
			}
		})
	}
}
