package assistant

import (
	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/model"
)

type SaveMode string

const (
	SaveModeNever     = SaveMode("NEVER")
	SaveModeFsLibrary = SaveMode("FS_LIBRARY")
)

var (
	ErrorNoAssistantEnabled       = errors.New("No assistant configured or enabled")
	ErrorMultipleProvidersEnabled = errors.New("only one provider can be enabled at a time")
)

// ProviderConfig represents a single LLM provider configuration.
type ProviderConfig struct {
	Type      ProviderType `yaml:"type" mapstructure:"type" head_comment:"Provider type: openai, anthropic, gemini, ollama, openai-compatible"`
	Enabled   bool         `yaml:"enabled" mapstructure:"enabled" head_comment:"If set to false, this provider will be skipped."`
	Model     string       `yaml:"model" mapstructure:"model" head_comment:"Model name to use (e.g., gpt-4o, claude-sonnet-4-20250514, gemini-1.5-flash)"`
	APIKeyEnv string       `yaml:"apiKeyEnv,omitempty" mapstructure:"apiKeyEnv" head_comment:"Environment variable holding the API key."`
	Endpoint  string       `yaml:"endpoint,omitempty" mapstructure:"endpoint" head_comment:"Custom API endpoint (optional, uses provider default if empty)."`
	// Ollama-specific
	ServerURL string `yaml:"serverUrl,omitempty" mapstructure:"serverUrl" head_comment:"Ollama server URL (for ollama provider only)."`
}

// Config is the top-level assistant configuration.
type Config struct {
	SaveMode  SaveMode         `yaml:"saveMode" mapstructure:"saveMode" head_comment:"Defines if you want to save the snippets created by the assistant. Possible values: NEVER | FS_LIBRARY"`
	Providers []ProviderConfig `yaml:"providers,omitempty" mapstructure:"providers" head_comment:"List of LLM providers. The first enabled provider will be used."`
}

// ValidateConfig returns an error if the configuration is invalid.
func (c Config) ValidateConfig() error {
	enabledCount := 0
	for _, p := range c.Providers {
		if p.Enabled {
			enabledCount++
		}
	}
	if enabledCount > 1 {
		return ErrorMultipleProvidersEnabled
	}
	return nil
}

// GetActiveProvider returns the first enabled provider config.
// Returns an error if multiple providers are enabled or none are enabled.
func (c Config) GetActiveProvider() (*ProviderConfig, error) {
	if err := c.ValidateConfig(); err != nil {
		return nil, err
	}
	for i := range c.Providers {
		if c.Providers[i].Enabled {
			return &c.Providers[i], nil
		}
	}
	return nil, ErrorNoAssistantEnabled
}

// ClientKey returns a key identifying the active provider.
func (c Config) ClientKey() (model.AssistantKey, error) {
	provider, err := c.GetActiveProvider()
	if err != nil {
		return "", err
	}
	return model.AssistantKey(provider.Type), nil
}

// HasEnabledProvider returns true if at least one provider is enabled.
func (c Config) HasEnabledProvider() bool {
	for _, p := range c.Providers {
		if p.Enabled {
			return true
		}
	}
	return false
}

// DefaultProviderConfig returns a default configuration for a given provider type.
func DefaultProviderConfig(providerType ProviderType) ProviderConfig {
	info := GetProviderInfo(providerType)
	cfg := ProviderConfig{
		Type:      providerType,
		Enabled:   true,
		Model:     info.DefaultModel,
		APIKeyEnv: info.DefaultAPIKeyEnv,
	}

	// Set Ollama-specific defaults
	if providerType == ProviderTypeOllama {
		cfg.ServerURL = "http://localhost:11434"
	}

	return cfg
}
