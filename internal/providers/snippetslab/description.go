package snippetslab

import "github.com/lemoony/snipkit/internal/model"

func Description(config Config) model.ProviderDescription {
	return model.ProviderDescription{
		Name:        "SnippetsLab",
		Description: "Use snippets form SnippetsLab",
		Enabled:     config.Enabled,
	}
}
