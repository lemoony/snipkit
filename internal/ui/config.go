package ui

import (
	"github.com/phuslu/log"
	"github.com/rivo/tview"
)

type Config struct {
	Theme        string       `yaml:"theme" head_comment:"The theme defines the terminal colors used by Snipkit.\nAvailable themes:default,dracula."`
	CustomThemes []NamedTheme `yaml:"customThemes" head_comment:"List of custom themes with values."`
}

type NamedTheme struct {
	Name   string      `yaml:"name"`
	Values ThemeValues `yaml:"values" head_comment:"A color can be created from a color name (W3C name) or by a hex value in the format #ffffff."`
}

type ThemeValues struct {
	BackgroundColor                       string
	BorderColor                           string
	BorderTitleColor                      string
	SyntaxHighlightingColorSchemeName     string
	SyntaxHighlightingApplyMainBackground bool
	PreviewSnippetDefaultTextColor        string
	ItemTextColor                         string
	SelectedItemTextColor                 string
	SelectedItemBackgroundColor           string
	ItemHighlightMatchColor               string
	CounterTextColor                      string
	CounterBackgroundColor                string
	LookupInputTextColor                  string
	LookupInputPlaceholderColor           string
	LookupInputBackgroundColor            string
	LookupLabelTextColor                  string
}

func (c *Config) GetSelectedTheme() ThemeValues {
	for i := range c.CustomThemes {
		if c.CustomThemes[i].Name == c.Theme {
			return c.CustomThemes[i].Values
		}
	}
	return themeDefault
}

var currentTheme ThemeValues

func DefaultConfig() Config {
	return Config{
		Theme: "default",
	}
}

func ApplyConfig(cfg Config) {
	theme := themeDefault
	for _, t := range cfg.CustomThemes {
		if t.Name == cfg.Theme {
			log.Debug().Msgf("Applied custom theme: %s", t.Name)
			theme = t.Values
		}
	}

	if t, ok := themeNames[cfg.Theme]; ok {
		log.Debug().Msgf("Applied provided theme: %s", cfg.Theme)
		theme = t
	}

	setTheme(theme)
}

func setTheme(theme ThemeValues) {
	currentTheme = theme
	tview.Styles.PrimitiveBackgroundColor = theme.backgroundColor()
	tview.Styles.BorderColor = theme.borderColor()
	tview.Styles.TitleColor = theme.borderTitleColor()
}
