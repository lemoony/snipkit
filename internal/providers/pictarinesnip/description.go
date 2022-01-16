package pictarinesnip

import "github.com/lemoony/snipkit/internal/model"

func Description(config Config) model.ProviderDescription {
	return model.ProviderDescription{
		Name:        "Pictarine Snip - Snippet Manager",
		Description: "Use snippets form Snip Snippets Manager (Pictarine)",
		Enabled:     config.Enabled,
	}
}
