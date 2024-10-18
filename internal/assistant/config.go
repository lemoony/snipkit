package assistant

import (
	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
)

type SaveMode string

const (
	SaveModeNever     = SaveMode("NEVER")
	SaveModeFsLibrary = SaveMode("FS_LIBRARY")
)

type Config struct {
	SaveMode SaveMode       `yaml:"saveMode" mapstructure:"saveMode" head_comment:"Defines if you want to save the snippets created by the assistant. Possible values: NEVER | FS_LIBRARY"`
	OpenAI   *openai.Config `yaml:"openai,omitempty" mapstructure:"openai"`
	Gemini   *gemini.Config `yaml:"gemini,omitempty" mapstructure:"gemini"`
}

func (c Config) moreThanOneEnabled() bool {
	openAIEnabled := c.OpenAI != nil && c.OpenAI.Enabled
	geminiEnabled := c.Gemini != nil && c.Gemini.Enabled
	return openAIEnabled && geminiEnabled
}
