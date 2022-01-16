package pictarinesnip

import "github.com/lemoony/snipkit/internal/model"

func Description(config Config) model.ManagerDescription {
	return model.ManagerDescription{
		Name:        "Pictarine Snip - Snippet Manager",
		Description: "Use snippets form Snip Snippets Manager (Pictarine)",
		Enabled:     config.Enabled,
	}
}
