package openai

type Config struct {
	Enabled   bool   `yaml:"enabled" head_comment:"If set to false, OpenAI will not be used as an AI assistant."`
	Endpoint  string `yaml:"endpoint" head_comment:"OpenAI API endpoint."`
	Model     string `yaml:"model" head_comment:"OpenAI Model to be used (e.g., openai/gpt-4o)"`
	APIKeyEnv string `yaml:"apiKeyEnv" head_comment:"The name of the environment variable holding the OpenAI API key."`
}

func AutoDiscoveryConfig(config *Config) Config {
	if config != nil {
		result := *config
		result.Enabled = true
		return result
	}

	return Config{
		Enabled:   false,
		Endpoint:  "https://api.openai.com",
		Model:     "openai/gpt-4o",
		APIKeyEnv: "SNIPKIT_OPENAPI_APIKEY",
	}
}
