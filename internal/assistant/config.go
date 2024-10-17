package assistant

import (
	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
)

type Config struct {
	OpenAI *openai.Config `yaml:"openai,omitempty" mapstructure:"openai"`
	Gemini *gemini.Config `yaml:"gemini,omitempty" mapstructure:"gemini"`
}
