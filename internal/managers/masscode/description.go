package masscode

import "github.com/lemoony/snipkit/internal/model"

const Key = model.ManagerKey("massCode")

func Description(config *Config) model.ManagerDescription {
	return model.ManagerDescription{
		Key:         Key,
		Name:        "massCode",
		Description: "Use Snippets form massCode - a free and open source code Snippets manager.",
		Enabled:     config != nil && config.Enabled,
	}
}
