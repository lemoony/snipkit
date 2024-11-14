package assistant

import (
	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/model"
)

type SaveMode string

const (
	SaveModeNever     = SaveMode("NEVER")
	SaveModeFsLibrary = SaveMode("FS_LIBRARY")

	noopKey = model.AssistantKey("noop")
)

var ErrorNoAssistantEnabled = errors.New("No assistant configured or enabled")

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

func (c Config) ClientKey() (model.AssistantKey, error) {
	switch {
	case c.moreThanOneEnabled():
		return noopKey, errors.New("Invalid config - more than one assistant is enabled")
	case c.OpenAI != nil && c.OpenAI.Enabled:
		return openai.Key, nil
	case c.Gemini != nil && c.Gemini.Enabled:
		return gemini.Key, nil
	default:
		return noopKey, ErrorNoAssistantEnabled
	}
}
