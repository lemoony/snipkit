package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	activeColor = lipgloss.Color("#F25D94")

	selectionColor   = lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}
	verySubduedColor = lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
)

type Style struct {
	width  int
	height int

	minimize bool

	needsToResize bool

	Reload func()
}

func DefaultStyle() Style {
	return Style{}
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

func (s *Style) SelectionColor() lipgloss.TerminalColor {
	return selectionColor
}

func (s *Style) SubduedColor() lipgloss.TerminalColor {
	return verySubduedColor
}

func (s *Style) ActiveColor() lipgloss.TerminalColor {
	return activeColor
}

func (s *Style) ButtonTextColor(selected bool) lipgloss.TerminalColor {
	if selected {
		return lipgloss.Color("#FFF7DB")
	}
	return lipgloss.Color("#FFF7DB")
}

func (s *Style) ButtonColor(selected bool) lipgloss.TerminalColor {
	if selected {
		return lipgloss.Color("#F25D94")
	}
	return lipgloss.Color("#888B7E")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
