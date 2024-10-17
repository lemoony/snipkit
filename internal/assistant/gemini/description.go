package gemini

import "github.com/lemoony/snipkit/internal/model"

const Key = model.AssistantKey("gemini")

func Description(config *Config) model.AssistantDescription {
	return model.AssistantDescription{
		Key:         Key,
		Name:        "Gemini",
		Description: "Use Google Gemini as an assistant AI",
		Enabled:     config != nil && config.Enabled,
	}
}
