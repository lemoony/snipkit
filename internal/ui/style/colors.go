package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
	"github.com/muesli/termenv"
)

var (
	colorProfile           = termenv.ColorProfile()
	defaultForegroundColor = lipgloss.Color(termenv.ConvertToRGB(termenv.ForegroundColor()).Hex())
	hasDarkBackground      = termenv.HasDarkBackground()
)

type Color struct {
	raw           string
	terminalColor lipgloss.TerminalColor
	tcellColor    tcell.Color
}

func newColor(val string) Color {
	return Color{
		raw:           val,
		terminalColor: lipgloss.Color(val),
		tcellColor:    tcell.GetColor(val),
	}
}

func (c Color) Value() lipgloss.TerminalColor {
	return c.terminalColor
}

func (c Color) CellValue() tcell.Color {
	return c.tcellColor
}

type colors struct {
	borderColor      Color
	borderTitleColor Color

	previewColorSchemeName string

	textColor Color

	subduedColor         Color
	subduedContrastColor Color

	verySubduedColor         Color
	verySubduedContrastColor Color

	activeColor         Color
	activeContrastColor Color

	titleColor         Color
	titleContrastColor Color

	highlightColor         Color
	highlightContrastColor Color

	infoColor         Color
	infoContrastColor Color

	snippetColor         Color
	snippetContrastColor Color

	successColor Color
	errorColor   Color
}

func newColors(t *ThemeValues) colors {
	return colors{
		borderColor:              newColor(t.BorderColor),
		borderTitleColor:         newColor(t.BorderTitleColor),
		previewColorSchemeName:   t.PreviewColorSchemeName,
		textColor:                colorWithDefaultForeground(t.TextColor),
		subduedColor:             newColor(t.SubduedColor),
		subduedContrastColor:     newColor(t.SubduedContrastColor),
		verySubduedColor:         newColor(t.VerySubduedColor),
		verySubduedContrastColor: newColor(t.VerySubduedContrastColor),
		activeColor:              newColor(t.ActiveColor),
		activeContrastColor:      newColor(t.ActiveContrastColor),
		titleColor:               newColor(t.TitleColor),
		titleContrastColor:       newColor(t.TitleContrastColor),
		highlightColor:           newColor(t.HighlightColor),
		infoColor:                newColor(t.InfoColor),
		infoContrastColor:        colorWithDefaultForeground(t.InfoContrastColor),
		snippetColor:             newColor(t.SnippetColor),
		snippetContrastColor:     colorWithDefaultForeground(t.SnippetContrastColor),
		successColor:             newColor(t.SuccessColor),
		errorColor:               newColor(t.ErrorColor),
	}
}

func colorWithDefaultForeground(val string) Color {
	if val == "" {
		return Color{
			terminalColor: defaultForegroundColor,
			tcellColor:    tcell.ColorReset,
		}
	}
	return newColor(val)
}
