package ui

import (
	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Config struct {
	Theme      string `yaml:"theme" head_comment:"The theme defines the terminal colors used by Snipkit.\nAvailable themes:default(.light|.dark),simple."`
	HideKeyMap bool   `yaml:"hideKeyMap,omitempty" head_comment:"If set to true, the key map won't be displayed. Default value: false"`
}

type NamedTheme struct {
	Name   string            `yaml:"name"`
	Values style.ThemeValues `yaml:"values" head_comment:"A color can be created from a color name (W3C name) or by a hex value in the format #ffffff."`
}

func (c *Config) GetSelectedTheme(system *system.System) style.ThemeValues {
	themeName := defaultThemeName
	if c.Theme != "" {
		themeName = c.Theme
	}

	if theme, ok := embeddedTheme(themeName); ok {
		return *theme
	}

	if theme, ok := customTheme(themeName, system); ok {
		return *theme
	}

	panic(errors.Wrapf(ErrInvalidTheme, "theme not found: %s", themeName))
}

func DefaultConfig() Config {
	return Config{Theme: "default"}
}

func ApplyConfig(cfg Config, system *system.System) {
	applyTheme(cfg, system)
}
