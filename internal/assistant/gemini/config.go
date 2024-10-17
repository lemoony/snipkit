package gemini

type Config struct {
	Enabled   bool   `yaml:"enabled" head_comment:"If set to false, Gemnini will not be used as an AI assistant."`
	Endpoint  string `yaml:"endpoint" head_comment:"Gemini API endpoint."`
	Model     string `yaml:"model" head_comment:"Gemini Model to be used (e.g., gemini-1.5-flash)"`
	APIKeyEnv string `yaml:"apiKeyEnv" head_comment:"The name of the environment variable holding the Gemini API key."`
}

func AutoDiscoveryConfig(config *Config) Config {
	if config != nil {
		result := *config
		result.Enabled = true
		return result
	}

	return Config{
		Enabled:   true,
		Endpoint:  "https://generativelanguage.googleapis.com",
		Model:     "gemini-1.5-flash",
		APIKeyEnv: "SNIPKIT_GEMINI_API_KEY",
	}
}
