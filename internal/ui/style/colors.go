package style

import "github.com/charmbracelet/lipgloss"

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

	selectionColor         lipgloss.TerminalColor
	selectionContrastColor lipgloss.TerminalColor

	titleColor         lipgloss.TerminalColor
	titleContrastColor lipgloss.TerminalColor

	highlightColor         lipgloss.TerminalColor
	highlightContrastColor lipgloss.TerminalColor

	infoColor         lipgloss.TerminalColor
	infoContrastColor lipgloss.TerminalColor
}

func newColors(t *ThemeValues) colors {
	return colors{
		borderColor:              lipgloss.Color(t.BorderColor),
		borderTitleColor:         lipgloss.Color(t.BorderTitleColor),
		previewColorSchemeName:   t.PreviewColorSchemeName,
		textColor:                lipgloss.Color(t.TextColor),
		subduedColor:             lipgloss.Color(t.SubduedColor),
		subduedContrastColor:     lipgloss.Color(t.SubduedContrastColor),
		verySubduedColor:         lipgloss.Color(t.VerySubduedColor),
		verySubduedContrastColor: lipgloss.Color(t.VerySubduedContrastColor),
		activeColor:              lipgloss.Color(t.ActiveColor),
		activeContrastColor:      lipgloss.Color(t.ActiveContrastColor),
		selectionColor:           lipgloss.Color(t.SelectionColor),
		selectionContrastColor:   lipgloss.Color(t.SelectionContrastColor),
		titleColor:               lipgloss.Color(t.TitleColor),
		titleContrastColor:       lipgloss.Color(t.TitleContrastColor),
		highlightColor:           lipgloss.Color(t.HighlightColor),
		infoColor:                lipgloss.Color(t.InfoColor),
		infoContrastColor:        lipgloss.Color(t.InfoContrastColor),
	}
}
