package ui

type Config struct {
	Theme        string       `yaml:"theme" head_comment:"The theme defines the terminal colors used by Snipkit.\nAvailable themes:dracula,solarized,light,dark."`
	CustomThemes []NamedTheme `yaml:"customThemes" head_comment:"List of custom themes with values."`
}

type NamedTheme struct {
	Name   string      `yaml:"name"`
	Values ThemeValues `yaml:"values" head_comment:"A color can be created from a color name (W3C name) or by a hex value in the format #ffffff."`
}

type ThemeValues struct {
	Background             string `yaml:"background"`
	ContrastBackground     string `yaml:"contrast_background"`
	MoreContrastBackground string `yaml:"moreContrastBackground"`
	Border                 string `yaml:"border"`
	Title                  string `yaml:"title"`
	Graphics               string `yaml:"graphics"`
	Text                   string `yaml:"text"`
	SecondaryText          string `yaml:"secondaryText"`
	TertiaryText           string `yaml:"tertiaryText"`
	InverseText            string `yaml:"inverseText"`
	ContrastSecondaryText  string `yaml:"contrastSecondaryText"`
}

func (c *Config) GetSelectedTheme() ThemeValues {
	for i := range c.CustomThemes {
		if c.CustomThemes[i].Name == c.Theme {
			return c.CustomThemes[i].Values
		}
	}

	return defaultThemeValues
}
