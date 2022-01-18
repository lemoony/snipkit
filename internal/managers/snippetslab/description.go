package snippetslab

import "github.com/lemoony/snipkit/internal/model"

var Key = model.ManagerKey("snippetslab")

func Description(config *Config) model.ManagerDescription {
	return model.ManagerDescription{
		Key:         Key,
		Name:        "SnippetsLab",
		Description: "Use snippets form SnippetsLab",
		Enabled:     config != nil && config.Enabled,
	}
}
