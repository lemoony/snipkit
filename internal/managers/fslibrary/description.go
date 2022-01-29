package fslibrary

import "github.com/lemoony/snipkit/internal/model"

const Key = model.ManagerKey("fslibrary")

func Description(config *Config) model.ManagerDescription {
	return model.ManagerDescription{
		Key:         Key,
		Name:        "File System Library",
		Description: "Use snippets form a local directory which holds snippet files",
		Enabled:     config != nil && config.Enabled,
	}
}
