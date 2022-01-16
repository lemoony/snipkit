package pictarinesnip

import "github.com/lemoony/snipkit/internal/model"

const Key = model.ManagerKey("pictarinesnip")

func Description(config *Config) model.ManagerDescription {
	return model.ManagerDescription{
		Key:         Key,
		Name:        "Pictarine Snip - Snippet Manager",
		Description: "Use snippets form Snip Snippets Manager (Pictarine)",
		Enabled:     config != nil && config.Enabled,
	}
}
