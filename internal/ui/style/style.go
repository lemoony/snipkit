package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type Style struct {
	width  int
	height int

	colors colors

	minimize bool

	needsToResize bool

	Reload func()
}

var NoopStyle = &Style{}

func NewStyle(t *ThemeValues) Style {
	return Style{
		colors: newColors(t),
	}
}

func (s *Style) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *Style) NeedsResize() bool {
	return s.needsToResize
}

func (s *Style) Title(text string) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		MarginBottom(1).
		Render(text)
}

func (s *Style) FormFieldWrapper(field string) string {
	if s.minimize {
		return lipgloss.NewStyle().Margin(0, 0, 0, 0).Render(field)
	}
	return lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(field)
}

func (s *Style) MainView(view string, help string, resize bool) string {
	defer func() {
		if resize {
			s.needsToResize = false
		}
	}()

	var sections []string
	sections = append(sections, view)

	marginsDefault := []int{1, 2, 1, 4}
	marginsMinimal := []int{0, 2, 0, 4}

	viewHeight := lipgloss.Height(view)
	helpHeight := lipgloss.Height(help)

	var margins []int
	if viewHeight+helpHeight+marginsDefault[0]+marginsDefault[2] <= s.height {
		margins = marginsDefault
		if s.minimize {
			s.minimize = false
			if !resize {
				s.needsToResize = true
			}
		}
	} else {
		margins = marginsMinimal
		if !s.minimize {
			s.minimize = true
			if !resize {
				s.needsToResize = true
			}
		}
	}

	// Fill empty space with newlines
	extraLines := ""
	if requiredHeight := viewHeight + helpHeight + margins[0] + margins[2]; requiredHeight < s.height {
		extraLines = strings.Repeat("\n", max(0, s.height-requiredHeight-1))
	}

	if extraLines != "" {
		sections = append(sections, extraLines)
	}

	sections = append(sections, help)

	return lipgloss.NewStyle().Margin(margins...).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (s *Style) ColorProfile() termenv.Profile {
	return lipgloss.ColorProfile()
}

func (s *Style) PreviewColorSchemeName() string {
	return s.colors.previewColorSchemeName
}

func (s *Style) BorderColor() lipgloss.TerminalColor {
	return s.colors.borderColor
}

func (s *Style) BorderTitleColor() lipgloss.TerminalColor {
	return s.colors.titleColor
}

func (s *Style) TextColor() lipgloss.TerminalColor {
	return s.colors.textColor
}

func (s *Style) PlaceholderColor() lipgloss.TerminalColor {
	return s.colors.subduedColor
}

func (s *Style) SelectionColor() lipgloss.TerminalColor {
	return s.colors.selectionColor
}

func (s *Style) SelectionColorReverse() lipgloss.TerminalColor {
	return s.colors.selectionContrastColor
}

func (s *Style) SubduedColor() lipgloss.TerminalColor {
	return s.colors.subduedColor
}

func (s *Style) ActiveColor() lipgloss.TerminalColor {
	return s.colors.activeColor
}

func (s *Style) ActiveContrastColor() lipgloss.TerminalColor {
	return s.colors.activeContrastColor
}

func (s *Style) InfoColor() lipgloss.TerminalColor {
	return s.colors.infoColor
}

func (s *Style) HighlightColor() lipgloss.TerminalColor {
	return s.colors.highlightColor
}

func (s *Style) HighlightContrastColor() lipgloss.TerminalColor {
	return s.colors.highlightContrastColor
}

func (s *Style) ButtonTextColor(selected bool) lipgloss.TerminalColor {
	if selected {
		return s.colors.activeContrastColor
	}
	return s.colors.subduedContrastColor
}

func (s *Style) ButtonColor(selected bool) lipgloss.TerminalColor {
	if selected {
		return s.colors.activeColor
	}
	return s.colors.subduedColor
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}