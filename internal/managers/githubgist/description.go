package githubgist

import "github.com/lemoony/snipkit/internal/model"

var Key = model.ManagerKey("githubgist")

func Description(config *Config) model.ManagerDescription {
	return model.ManagerDescription{
		Key:         Key,
		Name:        "Github Gist",
		Description: "Use snippets form Github Gist",
		Enabled:     config != nil && config.Enabled,
	}
}
