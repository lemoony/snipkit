package gemini

type Config struct {
	Enabled   bool   `yaml:"enabled" head_comment:"If set to false, Gemnini will not be used as an AI assistant."`
	Endpoint  string `yaml:"endpoint" head_comment:"Gemini API endpoint."`
	Version   string `yaml:"version" head_comment:"Gemini API version - currently only v1beta is supported."`
	Model     string `yaml:"model" head_comment:"Gemini Model to be used (e.g., gemini-1.5-flash)"`
	APIKeyEnv string `yaml:"apiKeyEnv" head_comment:"The name of the environment variable holding the Gemini API key."`
}

// AutoDiscoveryConfig is directly used by migrate so duplicate this function if the config changes.
func AutoDiscoveryConfig(config *Config) Config {
	if config != nil {
		result := *config
		result.Enabled = true
		return result
	}

	return Config{
		Enabled:   true,
		Endpoint:  "https://generativelanguage.googleapis.com",
		Version:   "v1beta",
		Model:     "gemini-1.5-flash",
		APIKeyEnv: "",
	}
}
