package openai

type Config struct {
	Enabled   bool   `yaml:"enabled" head_comment:"If set to false, OpenAI will not be used as an AI assistant."`
	Endpoint  string `yaml:"endpoint" head_comment:"OpenAI API endpoint."`
	Version   string `yaml:"version" head_comment:"OpenAI API version - currently ony v1 is supported."`
	Model     string `yaml:"model" head_comment:"OpenAI Model to be used (e.g., openai/gpt-4o)"`
	APIKeyEnv string `yaml:"apiKeyEnv" head_comment:"The name of the environment variable holding the OpenAI API key. If empty, SnipKit will ask for the API key and cache it."`
}

// AutoDiscoveryConfig is directly used by migrate so duplicate this function if the config changes.
func AutoDiscoveryConfig() *Config {
	return &Config{
		Enabled:   false,
		Endpoint:  "https://api.openai.com",
		Version:   "v1",
		Model:     "openai/gpt-4o",
		APIKeyEnv: "",
	}
}
