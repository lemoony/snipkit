package ui

import (
	"emperror.dev/errors"
	"github.com/rivo/tview"

	"github.com/lemoony/snippet-kit/internal/utils/system"
)

type Config struct {
	Theme string `yaml:"theme" head_comment:"The theme defines the terminal colors used by Snipkit.\nAvailable themes:default,dracula."`
}

type NamedTheme struct {
	Name   string      `yaml:"name"`
	Values ThemeValues `yaml:"values" head_comment:"A color can be created from a color name (W3C name) or by a hex value in the format #ffffff."`
}

type ThemeValues struct {
	BackgroundColor                              string `yaml:"backgroundColor"`
	BorderColor                                  string `yaml:"borderColor"`
	BorderTitleColor                             string `yaml:"borderTitleColor"`
	PreviewColorSchemeName                       string `yaml:"previewColorSchemeName"`
	PreviewApplyMainBackground                   bool   `yaml:"previewApplyMainBackground"`
	PreviewOverwriteBackgroundColor              string `yaml:"previewOverwriteBackgroundColor"`
	PreviewDefaultTextColor                      string `yaml:"previewDefaultTextColor"`
	ItemTextColor                                string `yaml:"itemTextColor"`
	SelectedItemTextColor                        string `yaml:"selectedItemTextColor"`
	SelectedItemBackgroundColor                  string `yaml:"selectedItemBackgroundColor"`
	ItemHighlightMatchBackgroundColor            string `yaml:"itemHighlightMatchBackgroundColor"`
	ItemHighlightMatchTextColor                  string `yaml:"itemHighlightMatchTextColor"`
	CounterTextColor                             string `yaml:"counterTextColor"`
	LookupInputTextColor                         string `yaml:"lookupInputTextColor"`
	LookupInputPlaceholderColor                  string `yaml:"lookupInputPlaceholderColor"`
	LookupInputBackgroundColor                   string `yaml:"lookupInputBackgroundColor"`
	ParametersLabelTextColor                     string `yaml:"parametersLabelTextColor"`
	ParametersFieldBackgroundColor               string `yaml:"parametersFieldBackgroundColor"`
	ParametersFieldTextColor                     string `yaml:"parametersFieldTextColor"`
	ParameterAutocompleteBackgroundColor         string `yaml:"parameterAutocompleteBackgroundColor"`
	ParameterAutocompleteTextColor               string `yaml:"parameterAutocompleteTextColor"`
	ParameterAutocompleteSelectedBackgroundColor string `yaml:"parameterAutocompleteSelectedBackgroundColor"`
	ParameterAutocompleteSelectedTextColor       string `yaml:"parameterAutocompleteSelectedTextColor"`
	SelectedButtonBackgroundColor                string `yaml:"selectedButtonBackgroundColor"`
	SelectedButtonTextColor                      string `yaml:"selectedButtonTextColor"`
}

func (c *Config) GetSelectedTheme(system *system.System) ThemeValues {
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

var currentTheme ThemeValues

func DefaultConfig() Config {
	return Config{
		Theme: "default",
	}
}

func ApplyConfig(cfg Config, system *system.System) {
	setTheme(cfg.GetSelectedTheme(system))
}

func setTheme(theme ThemeValues) {
	currentTheme = theme
	tview.Styles.PrimitiveBackgroundColor = theme.backgroundColor()
	tview.Styles.BorderColor = theme.borderColor()
	tview.Styles.TitleColor = theme.borderTitleColor()
}
