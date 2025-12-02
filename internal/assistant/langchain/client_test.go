package langchain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"

	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
)

func Test_extractAPIError_NilAndSimple(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		result := extractAPIError(nil)
		assert.Nil(t, result)
	})

	t.Run("simple error", func(t *testing.T) {
		err := errors.New("simple error")
		result := extractAPIError(err)
		assert.Equal(t, "simple error", result.Error())
	})
}

func Test_extractAPIError_LLMsError(t *testing.T) {
	tests := []struct {
		name     string
		err      *llms.Error
		expected string
	}{
		{
			name:     "with message",
			err:      &llms.Error{Code: llms.ErrCodeAuthentication, Message: "invalid API key", Provider: "openai"},
			expected: "openai: invalid API key",
		},
		{
			name: "with details",
			err: &llms.Error{
				Code: llms.ErrCodeRateLimit, Message: "rate limit exceeded", Provider: "anthropic",
				Details: map[string]interface{}{"retry_after": 60},
			},
			expected: "anthropic: rate limit exceeded (details: map[retry_after:60])",
		},
		{
			name:     "with cause",
			err:      &llms.Error{Code: llms.ErrCodeInvalidRequest, Message: "bad request", Provider: "gemini", Cause: errors.New("underlying error")},
			expected: "gemini: bad request (cause: underlying error)",
		},
		{
			name:     "without message uses code",
			err:      &llms.Error{Code: llms.ErrCodeTimeout, Provider: "ollama"},
			expected: "ollama: timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractAPIError(tt.err)
			assert.Equal(t, tt.expected, result.Error())
		})
	}
}

func Test_extractAPIError_ReflectionBased(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "Code and Message fields",
			err:      &errorWithCodeAndMessage{Code: 404, Message: "model not found"},
			expected: "API error 404: model not found",
		},
		{
			name:     "Code and Body fields",
			err:      &errorWithCodeAndBody{Code: 500, Body: `{"error": "internal server error"}`},
			expected: `API error 500: {"error": "internal server error"}`,
		},
		{
			name:     "only Message field",
			err:      &errorWithMessage{Message: "something went wrong"},
			expected: "API error: something went wrong",
		},
		{
			name:     "only Body field",
			err:      &errorWithBody{Body: "raw error body"},
			expected: "API error: raw error body",
		},
		{
			name:     "wrapped error with details",
			err:      wrapError(&errorWithCodeAndMessage{Code: 403, Message: "forbidden"}),
			expected: "API error 403: forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractAPIError(tt.err)
			assert.Equal(t, tt.expected, result.Error())
		})
	}
}

// Test error types for reflection-based extraction

type errorWithCodeAndMessage struct {
	Code    int
	Message string
}

func (e *errorWithCodeAndMessage) Error() string {
	return "error with code and message"
}

type errorWithCodeAndBody struct {
	Code int
	Body string
}

func (e *errorWithCodeAndBody) Error() string {
	return "error with code and body"
}

type errorWithMessage struct {
	Message string
}

func (e *errorWithMessage) Error() string {
	return "error with message"
}

type errorWithBody struct {
	Body string
}

func (e *errorWithBody) Error() string {
	return "error with body"
}

type wrappedError struct {
	cause error
}

func (e *wrappedError) Error() string {
	return "wrapped: " + e.cause.Error()
}

func (e *wrappedError) Unwrap() error {
	return e.cause
}

func wrapError(err error) error {
	return &wrappedError{cause: err}
}

// mockLLMModel is a mock implementation of llms.Model for testing.
type mockLLMModel struct {
	responseContent string
	err             error
	emptyResponse   bool
}

func (m *mockLLMModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	if m.err != nil {
		return nil, m.err
	}

	if m.emptyResponse {
		return &llms.ContentResponse{Choices: []*llms.ContentChoice{}}, nil
	}

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{Content: m.responseContent},
		},
	}, nil
}

func (m *mockLLMModel) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return "", nil
}

func Test_getAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		envValue    string
		setEnv      bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty env var name",
			envVar:      "",
			expectError: true,
			errorMsg:    "no API key environment variable specified",
		},
		{
			name:        "env var not set",
			envVar:      "NONEXISTENT_API_KEY",
			expectError: true,
			errorMsg:    "environment variable NONEXISTENT_API_KEY is not set or empty",
		},
		{
			name:        "env var set to empty string",
			envVar:      "EMPTY_KEY",
			envValue:    "",
			setEnv:      true,
			expectError: true,
			errorMsg:    "environment variable EMPTY_KEY is not set or empty",
		},
		{
			name:        "valid API key",
			envVar:      "VALID_KEY",
			envValue:    "sk-test-123456",
			setEnv:      true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				cleanup := mockutil.MockAPIKeyEnv(tt.envVar, tt.envValue)
				defer cleanup()
			}

			result, err := getAPIKey(tt.envVar)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.envValue, result)
			}
		})
	}
}

//nolint:funlen // test case table
func getNewClientTestCases() []struct {
	name        string
	config      Config
	setupEnv    func() func() // Returns cleanup function
	expectError bool
	errorMsg    string
} {
	return []struct {
		name        string
		config      Config
		setupEnv    func() func() // Returns cleanup function
		expectError bool
		errorMsg    string
	}{
		{
			name: "openai - success",
			config: Config{
				Type:      ProviderTypeOpenAI,
				Model:     "gpt-4o",
				APIKeyEnv: "TEST_OPENAI_KEY",
			},
			setupEnv: func() func() {
				return mockutil.MockAPIKeyEnv("TEST_OPENAI_KEY", "sk-test-123")
			},
			expectError: false,
		},
		{
			name: "openai - missing API key",
			config: Config{
				Type:      ProviderTypeOpenAI,
				Model:     "gpt-4o",
				APIKeyEnv: "MISSING_KEY",
			},
			expectError: true,
			errorMsg:    "environment variable MISSING_KEY is not set",
		},
		{
			name: "openai - with custom endpoint",
			config: Config{
				Type:      ProviderTypeOpenAI,
				Model:     "gpt-4o",
				APIKeyEnv: "TEST_OPENAI_KEY",
				Endpoint:  "https://custom.openai.com/v1",
			},
			setupEnv: func() func() {
				return mockutil.MockAPIKeyEnv("TEST_OPENAI_KEY", "sk-test-123")
			},
			expectError: false,
		},
		{
			name: "anthropic - success",
			config: Config{
				Type:      ProviderTypeAnthropic,
				Model:     "claude-sonnet-4.5",
				APIKeyEnv: "TEST_ANTHROPIC_KEY",
			},
			setupEnv: func() func() {
				return mockutil.MockAPIKeyEnv("TEST_ANTHROPIC_KEY", "sk-ant-test")
			},
			expectError: false,
		},
		{
			name: "anthropic - with custom endpoint",
			config: Config{
				Type:      ProviderTypeAnthropic,
				Model:     "claude-sonnet-4.5",
				APIKeyEnv: "TEST_ANTHROPIC_KEY",
				Endpoint:  "https://custom.anthropic.com",
			},
			setupEnv: func() func() {
				return mockutil.MockAPIKeyEnv("TEST_ANTHROPIC_KEY", "sk-ant-test")
			},
			expectError: false,
		},
		{
			name: "gemini - success",
			config: Config{
				Type:      ProviderTypeGemini,
				Model:     "gemini-1.5-flash",
				APIKeyEnv: "TEST_GEMINI_KEY",
			},
			setupEnv: func() func() {
				return mockutil.MockAPIKeyEnv("TEST_GEMINI_KEY", "test-gemini-key")
			},
			expectError: false,
		},
		{
			name: "ollama - success (no API key required)",
			config: Config{
				Type:      ProviderTypeOllama,
				Model:     "llama3",
				ServerURL: "http://localhost:11434",
			},
			expectError: false,
		},
		{
			name: "ollama - with default server URL",
			config: Config{
				Type:  ProviderTypeOllama,
				Model: "llama3",
			},
			expectError: false,
		},
		{
			name: "openai-compatible - with API key",
			config: Config{
				Type:      ProviderTypeOpenAICompatible,
				Model:     "custom-model",
				APIKeyEnv: "TEST_COMPATIBLE_KEY",
				Endpoint:  "https://openrouter.ai/api/v1",
			},
			setupEnv: func() func() {
				return mockutil.MockAPIKeyEnv("TEST_COMPATIBLE_KEY", "sk-or-test")
			},
			expectError: false,
		},
		{
			name: "openai-compatible - without API key (currently fails - SDK limitation)",
			config: Config{
				Type:     ProviderTypeOpenAICompatible,
				Model:    "custom-model",
				Endpoint: "https://local-llm.com/v1",
			},
			// Note: The openai.New() SDK requires a token even when we don't provide one
			// This is a limitation of the underlying library
			expectError: true,
			errorMsg:    "missing the OpenAI API key",
		},
		{
			name: "openai-compatible - missing endpoint",
			config: Config{
				Type:      ProviderTypeOpenAICompatible,
				Model:     "custom-model",
				APIKeyEnv: "TEST_KEY",
			},
			setupEnv: func() func() {
				return mockutil.MockAPIKeyEnv("TEST_KEY", "sk-test")
			},
			expectError: true,
			errorMsg:    "endpoint is required for openai-compatible provider",
		},
		{
			name: "unsupported provider type",
			config: Config{
				Type:  ProviderType("unknown"),
				Model: "test",
			},
			expectError: true,
			errorMsg:    "unsupported provider type: unknown",
		},
	}
}

func Test_NewClient_AllProviders(t *testing.T) {
	tests := getNewClientTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				cleanup := tt.setupEnv()
				defer cleanup()
			}

			client, err := NewClient(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.model)
				assert.Empty(t, client.history) // History should be empty initially
			}
		})
	}
}

func testFirstQueryAddsSystemPrompt(t *testing.T) {
	mockModel := &mockLLMModel{
		responseContent: "Test response",
	}

	client := &Client{
		model:   mockModel,
		history: []llms.MessageContent{},
	}

	response, err := client.Query("test prompt")

	assert.NoError(t, err)
	assert.Equal(t, "Test response", response)

	// Verify history: system prompt + user message + AI response
	assert.Len(t, client.history, 3)
	assert.Equal(t, llms.ChatMessageTypeSystem, client.history[0].Role)
	assert.Equal(t, llms.ChatMessageTypeHuman, client.history[1].Role)
	assert.Equal(t, llms.ChatMessageTypeAI, client.history[2].Role)
}

func testSubsequentQueriesMaintainHistory(t *testing.T) {
	mockModel := &mockLLMModel{
		responseContent: "Response 2",
	}

	client := &Client{
		model: mockModel,
		history: []llms.MessageContent{
			{Role: llms.ChatMessageTypeSystem, Parts: []llms.ContentPart{llms.TextPart("system")}},
			{Role: llms.ChatMessageTypeHuman, Parts: []llms.ContentPart{llms.TextPart("first")}},
			{Role: llms.ChatMessageTypeAI, Parts: []llms.ContentPart{llms.TextPart("Response 1")}},
		},
	}

	response, err := client.Query("second prompt")

	assert.NoError(t, err)
	assert.Equal(t, "Response 2", response)
	assert.Len(t, client.history, 5)

	// Verify the new messages
	humanMsg := client.history[3]
	assert.Equal(t, llms.ChatMessageTypeHuman, humanMsg.Role)
	assert.Len(t, humanMsg.Parts, 1)
	textContent, ok := humanMsg.Parts[0].(llms.TextContent)
	assert.True(t, ok)
	assert.Equal(t, "second prompt", textContent.Text)

	aiMsg := client.history[4]
	assert.Equal(t, llms.ChatMessageTypeAI, aiMsg.Role)
	assert.Len(t, aiMsg.Parts, 1)
	textContent, ok = aiMsg.Parts[0].(llms.TextContent)
	assert.True(t, ok)
	assert.Equal(t, "Response 2", textContent.Text)
}

func Test_Client_Query_HistoryManagement(t *testing.T) {
	t.Run("first query adds system prompt", testFirstQueryAddsSystemPrompt)
	t.Run("subsequent queries maintain history", testSubsequentQueriesMaintainHistory)

	t.Run("API error handling", func(t *testing.T) {
		mockModel := &mockLLMModel{
			err: &llms.Error{
				Code:     llms.ErrCodeAuthentication,
				Message:  "invalid API key",
				Provider: "openai",
			},
		}

		client := &Client{model: mockModel, history: []llms.MessageContent{}}

		response, err := client.Query("test")

		assert.Error(t, err)
		assert.Empty(t, response)
		assert.Contains(t, err.Error(), "openai: invalid API key")
	})

	t.Run("no response choices", func(t *testing.T) {
		mockModel := &mockLLMModel{
			emptyResponse: true,
		}

		client := &Client{model: mockModel, history: []llms.MessageContent{}}

		response, err := client.Query("test")

		assert.Error(t, err)
		assert.Empty(t, response)
		assert.Contains(t, err.Error(), "no response from LLM")
	})
}
