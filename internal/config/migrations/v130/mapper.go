package config

import (
	"emperror.dev/errors"
	"gopkg.in/yaml.v3"
)

const (
	VersionFrom = "1.2.0"
	VersionTo   = "1.3.0"
)

func Migrate(old []byte) []byte {
	var config versionWrapper

	if err := yaml.Unmarshal(old, &config); err != nil {
		panic(err)
	}

	if config.Version != VersionFrom {
		panic(errors.Errorf("Invalid version for migration to v1.3.0: %s", config.Version))
	}

	config.Version = VersionTo

	// Convert old assistant config to new providers array format
	providers := []providerConfig{}

	if openai, ok := config.Config.Assistant["openai"].(map[string]interface{}); ok {
		provider := convertLegacyProvider("openai", openai)
		providers = append(providers, provider)
	}

	if gemini, ok := config.Config.Assistant["gemini"].(map[string]interface{}); ok {
		provider := convertLegacyProvider("gemini", gemini)
		providers = append(providers, provider)
	}

	// Preserve saveMode, set new providers array
	newAssistant := map[string]interface{}{
		"saveMode": config.Config.Assistant["saveMode"],
	}

	if len(providers) > 0 {
		newAssistant["providers"] = providers
	}

	config.Config.Assistant = newAssistant

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	return configBytes
}

func convertLegacyProvider(providerType string, legacy map[string]interface{}) providerConfig {
	return providerConfig{
		Type:      providerType,
		Enabled:   getBool(legacy, "enabled"),
		Model:     getString(legacy, "model"),
		APIKeyEnv: getString(legacy, "apiKeyEnv"),
		Endpoint:  getString(legacy, "endpoint"),
	}
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

type providerConfig struct {
	Type      string `yaml:"type"`
	Enabled   bool   `yaml:"enabled"`
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"apiKeyEnv,omitempty"`
	Endpoint  string `yaml:"endpoint,omitempty"`
}

type versionWrapper struct {
	Version string     `yaml:"version"`
	Config  configV130 `yaml:"config"`
}

type configV130 struct {
	Style              map[string]interface{} `yaml:"style"`
	Editor             string                 `yaml:"editor"`
	DefaultRootCommand string                 `yaml:"defaultRootCommand"`
	FuzzySearch        bool                   `yaml:"fuzzySearch"`
	SecretStorage      string                 `yaml:"secretStorage"`
	Script             map[string]interface{} `yaml:"scripts"`
	Assistant          map[string]interface{} `yaml:"assistant"`
	Manager            map[string]interface{} `yaml:"manager"`
}
