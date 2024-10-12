package ai

import (
	"github.com/lemoony/snipkit/internal/ai/openai"
)

type Config struct {
	OpenAI openai.Config `yaml:"openai,omitempty" mapstructure:"openai"`
}
