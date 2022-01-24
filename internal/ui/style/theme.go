package style

type ThemeValues struct {
	BorderColor      string `yaml:"borderColor"`
	BorderTitleColor string `yaml:"borderTitleColor"`

	PreviewColorSchemeName string `yaml:"previewColorSchemeName"`

	TextColor string `yaml:"textColor"`

	SubduedColor         string `yaml:"subduedColor"`
	SubduedContrastColor string `yaml:"subduedContrastColor"`

	VerySubduedColor         string `yaml:"verySubduedColor"`
	VerySubduedContrastColor string `yaml:"verySubduedContrastColor"`

	ActiveColor         string `yaml:"activeColor"`
	ActiveContrastColor string `yaml:"activeContrastColor"`

	TitleColor         string `yaml:"titleColor"`
	TitleContrastColor string `yaml:"titleContrastColor"`

	HighlightColor         string `yaml:"highlightColor"`
	HighlightContrastColor string `yaml:"highlightContrastColor"`

	InfoColor         string `yaml:"infoColor"`
	InfoContrastColor string `yaml:"infoContrastColor"`

	SnippetColor         string `yaml:"snippetColor"`
	SnippetContrastColor string `yaml:"snippetContrastColor"`
}
