package pet

import "github.com/lemoony/snipkit/internal/model"

const Key = model.ManagerKey("pet")

func Description(config *Config) model.ManagerDescription {
	return model.ManagerDescription{
		Key:         Key,
		Name:        "pet - CLI Snippet Manager",
		Description: "Use snippets form pet - a simple command-line snippet manager",
		Enabled:     config != nil && config.Enabled,
	}
}
