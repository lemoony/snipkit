package fslibrary

import "github.com/lemoony/snipkit/internal/model"

func Description(config *Config) model.ManagerDescription {
	return model.ManagerDescription{
		Name:        "File System Library",
		Description: "Use snippets form a local directory which holds snippet files",
		Enabled:     config != nil && config.Enabled,
	}
}
