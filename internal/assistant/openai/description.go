package openai

import "github.com/lemoony/snipkit/internal/model"

const Key = model.AssistantKey("openai")

func Description(config *Config) model.AssistantDescription {
	return model.AssistantDescription{
		Key:         Key,
		Name:        "OpenAI",
		Description: "Use OpenAI as an assistant AI",
		Enabled:     config != nil && config.Enabled,
	}
}
