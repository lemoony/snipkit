package assistant

import (
	"github.com/lemoony/snipkit/internal/assistant/openai"
)

type Config struct {
	OpenAI openai.Config `yaml:"openai,omitempty" mapstructure:"openai"`
}
