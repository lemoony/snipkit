package langchain

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/lemoony/snipkit/internal/assistant/prompts"
)

// ProviderType represents the type of LLM provider.
type ProviderType string

const (
	ProviderTypeOpenAI           = ProviderType("openai")
	ProviderTypeAnthropic        = ProviderType("anthropic")
	ProviderTypeGemini           = ProviderType("gemini")
	ProviderTypeOllama           = ProviderType("ollama")
	ProviderTypeOpenAICompatible = ProviderType("openai-compatible")
)

// Config holds the configuration for creating a langchain client.
type Config struct {
	Type      ProviderType
	Model     string
	APIKeyEnv string
	Endpoint  string
	ServerURL string // Ollama-specific
}

// Client wraps a langchaingo LLM model with conversation history support.
type Client struct {
	model   llms.Model
	history []llms.MessageContent
}

// NewClient creates a new langchain-based client for the given configuration.
func NewClient(cfg Config) (*Client, error) {
	model, err := createModel(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{model: model}, nil
}

func createModel(cfg Config) (llms.Model, error) {
	switch cfg.Type {
	case ProviderTypeOpenAI:
		return createOpenAIModel(cfg)
	case ProviderTypeAnthropic:
		return createAnthropicModel(cfg)
	case ProviderTypeGemini:
		return createGeminiModel(cfg)
	case ProviderTypeOllama:
		return createOllamaModel(cfg)
	case ProviderTypeOpenAICompatible:
		return createOpenAICompatibleModel(cfg)
	default:
		return nil, errors.Errorf("unsupported provider type: %s", cfg.Type)
	}
}

func createOpenAIModel(cfg Config) (llms.Model, error) {
	apiKey, err := getAPIKey(cfg.APIKeyEnv)
	if err != nil {
		return nil, err
	}

	opts := []openai.Option{
		openai.WithToken(apiKey),
		openai.WithModel(cfg.Model),
	}

	if cfg.Endpoint != "" {
		opts = append(opts, openai.WithBaseURL(cfg.Endpoint))
	}

	return openai.New(opts...)
}

func createAnthropicModel(cfg Config) (llms.Model, error) {
	apiKey, err := getAPIKey(cfg.APIKeyEnv)
	if err != nil {
		return nil, err
	}

	opts := []anthropic.Option{
		anthropic.WithToken(apiKey),
		anthropic.WithModel(cfg.Model),
	}

	if cfg.Endpoint != "" {
		opts = append(opts, anthropic.WithBaseURL(cfg.Endpoint))
	}

	return anthropic.New(opts...)
}

func createGeminiModel(cfg Config) (llms.Model, error) {
	apiKey, err := getAPIKey(cfg.APIKeyEnv)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	opts := []googleai.Option{
		googleai.WithAPIKey(apiKey),
		googleai.WithDefaultModel(cfg.Model),
	}

	return googleai.New(ctx, opts...)
}

func createOllamaModel(cfg Config) (llms.Model, error) {
	opts := []ollama.Option{
		ollama.WithModel(cfg.Model),
	}

	if cfg.ServerURL != "" {
		opts = append(opts, ollama.WithServerURL(cfg.ServerURL))
	}

	return ollama.New(opts...)
}

func createOpenAICompatibleModel(cfg Config) (llms.Model, error) {
	opts := []openai.Option{
		openai.WithModel(cfg.Model),
	}

	// API key is optional for some OpenAI-compatible endpoints
	if cfg.APIKeyEnv != "" {
		apiKey, err := getAPIKey(cfg.APIKeyEnv)
		if err != nil {
			return nil, err
		}
		opts = append(opts, openai.WithToken(apiKey))
	}

	if cfg.Endpoint != "" {
		opts = append(opts, openai.WithBaseURL(cfg.Endpoint))
	} else {
		return nil, errors.New("endpoint is required for openai-compatible provider")
	}

	return openai.New(opts...)
}

// Query sends a prompt to the LLM and returns the response.
// It maintains conversation history for multi-turn conversations.
func (c *Client) Query(prompt string) (string, error) {
	ctx := context.Background()
	log.Debug().Str("prompt", prompt).Msg("Starting query request")

	// Initialize with system prompt on first query
	if len(c.history) == 0 {
		c.history = append(c.history, llms.MessageContent{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextPart(prompts.DefaultPrompt)},
		})
	}

	// Add user message
	c.history = append(c.history, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextPart(prompt)},
	})

	log.Trace().Int("history_length", len(c.history)).Msg("Sending request to LLM")

	// Generate response
	response, err := c.model.GenerateContent(ctx, c.history)
	if err != nil {
		return "", extractAPIError(err)
	}

	if len(response.Choices) == 0 {
		return "", errors.New("no response from LLM")
	}

	aiResponse := response.Choices[0].Content
	log.Trace().Str("response", aiResponse).Msg("Received response from LLM")

	// Add AI response to history for multi-turn conversation
	c.history = append(c.history, llms.MessageContent{
		Role:  llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{llms.TextPart(aiResponse)},
	})

	return aiResponse, nil
}

func getAPIKey(envVar string) (string, error) {
	if envVar == "" {
		return "", errors.New("no API key environment variable specified in config")
	}
	apiKey := os.Getenv(envVar)
	if apiKey == "" {
		return "", errors.Errorf("environment variable %s is not set or empty - please set it with your API key", envVar)
	}
	return apiKey, nil
}

// extractAPIError extracts detailed error information from API errors.
// Uses reflection to extract common error fields (Message, Body, Details) in a provider-agnostic way.
func extractAPIError(err error) error {
	if err == nil {
		return nil
	}

	// Check for langchaingo's standardized error type first
	var llmErr *llms.Error
	if errors.As(err, &llmErr) {
		msg := llmErr.Message
		if msg == "" {
			msg = string(llmErr.Code)
		}
		if len(llmErr.Details) > 0 {
			return fmt.Errorf("%s: %s (details: %v)", llmErr.Provider, msg, llmErr.Details)
		}
		if llmErr.Cause != nil {
			return fmt.Errorf("%s: %s (cause: %v)", llmErr.Provider, msg, llmErr.Cause)
		}
		return fmt.Errorf("%s: %s", llmErr.Provider, msg)
	}

	// Use reflection to extract common error fields from any error type
	details := extractErrorDetails(err)
	if details != "" {
		return fmt.Errorf("%s", details)
	}

	return err
}

// extractErrorDetails uses reflection to find common error fields like Message, Body, Details.
func extractErrorDetails(err error) string {
	// Traverse the error chain
	for unwrapped := err; unwrapped != nil; unwrapped = errors.Unwrap(unwrapped) {
		if details := extractFieldsFromError(unwrapped); details != "" {
			return details
		}
	}
	return ""
}

// extractFieldsFromError extracts Message or Body fields from an error using reflection.
func extractFieldsFromError(err error) string {
	v := reflect.ValueOf(err)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	// Try to get common error fields
	var code, message, body string

	if f := v.FieldByName("Code"); f.IsValid() && f.CanInterface() {
		code = fmt.Sprintf("%v", f.Interface())
	}
	if f := v.FieldByName("Message"); f.IsValid() && f.Kind() == reflect.String {
		message = f.String()
	}
	if f := v.FieldByName("Body"); f.IsValid() && f.Kind() == reflect.String {
		body = f.String()
	}

	// Build detailed error message
	if message != "" {
		if code != "" {
			return fmt.Sprintf("API error %s: %s", code, message)
		}
		return fmt.Sprintf("API error: %s", message)
	}
	if body != "" {
		if code != "" {
			return fmt.Sprintf("API error %s: %s", code, body)
		}
		return fmt.Sprintf("API error: %s", body)
	}

	return ""
}
