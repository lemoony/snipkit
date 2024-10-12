package openai

import (
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Config struct {
	Enabled  bool   `yaml:"enabled" head_comment:"If set to false, the files specified via libraryPath will not be provided to you."`
	Endpoint string `yaml:"endpoint" head_comment:"OpenAI API endpoint."`
	Version  string `yaml:"version" head_comment:"OpenAI API version - currently ony v1 is supported."`
	Model    string `yaml:"model" head_comment:"OpenAI Model to be used (e.g., openai/gpt-4o)"`
}

func AutoDiscoveryConfig(system *system.System) *Config {
	return &Config{
		Enabled:  false,
		Endpoint: "https://api.openai.com",
		Version:  "v1",
		Model:    "openai/gpt-4o",
	}
}
