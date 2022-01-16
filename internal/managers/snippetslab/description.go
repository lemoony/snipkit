package snippetslab

import "github.com/lemoony/snipkit/internal/model"

func Description(config Config) model.ManagerDescription {
	return model.ManagerDescription{
		Name:        "SnippetsLab",
		Description: "Use snippets form SnippetsLab",
		Enabled:     config.Enabled,
	}
}
