package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	colorProfile           = termenv.ColorProfile()
	defaultForegroundColor = lipgloss.Color(termenv.ConvertToRGB(termenv.ForegroundColor()).Hex())
	hasDarkBackground      = termenv.HasDarkBackground()
)

type colors struct {
	borderColor      lipgloss.TerminalColor
	borderTitleColor lipgloss.TerminalColor

	previewColorSchemeName string

	textColor lipgloss.TerminalColor

	subduedColor         lipgloss.TerminalColor
	subduedContrastColor lipgloss.TerminalColor

	verySubduedColor         lipgloss.TerminalColor
	verySubduedContrastColor lipgloss.TerminalColor

	activeColor         lipgloss.TerminalColor
	activeContrastColor lipgloss.TerminalColor

	titleColor         lipgloss.TerminalColor
	titleContrastColor lipgloss.TerminalColor

	highlightColor         lipgloss.TerminalColor
	highlightContrastColor lipgloss.TerminalColor

	infoColor         lipgloss.TerminalColor
	infoContrastColor lipgloss.TerminalColor

	snippetColor         lipgloss.TerminalColor
	snippetContrastColor lipgloss.TerminalColor
}

func newColors(t *ThemeValues) colors {
	return colors{
		borderColor:              lipgloss.Color(t.BorderColor),
		borderTitleColor:         lipgloss.Color(t.BorderTitleColor),
		previewColorSchemeName:   t.PreviewColorSchemeName,
		textColor:                color(t.TextColor, defaultForegroundColor),
		subduedColor:             lipgloss.Color(t.SubduedColor),
		subduedContrastColor:     lipgloss.Color(t.SubduedContrastColor),
		verySubduedColor:         lipgloss.Color(t.VerySubduedColor),
		verySubduedContrastColor: lipgloss.Color(t.VerySubduedContrastColor),
		activeColor:              lipgloss.Color(t.ActiveColor),
		activeContrastColor:      lipgloss.Color(t.ActiveContrastColor),
		titleColor:               lipgloss.Color(t.TitleColor),
		titleContrastColor:       lipgloss.Color(t.TitleContrastColor),
		highlightColor:           lipgloss.Color(t.HighlightColor),
		infoColor:                lipgloss.Color(t.InfoColor),
		infoContrastColor:        color(t.InfoContrastColor, defaultForegroundColor),
		snippetColor:             lipgloss.Color(t.SnippetColor),
		snippetContrastColor:     color(t.SnippetContrastColor, defaultForegroundColor),
	}
}

func color(val string, defaultColor lipgloss.Color) lipgloss.Color {
	if val == "" {
		return defaultColor
	}
	return lipgloss.Color(val)
}
