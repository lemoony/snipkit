package langchain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
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
