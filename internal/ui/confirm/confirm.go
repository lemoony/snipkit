package confirm

import (
	"bytes"
	"io"
	"text/template"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"github.com/muesli/termenv"

	"github.com/lemoony/snipkit/internal/ui/uimsg"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2, 0, 4) //nolint:gomnd // magic number is okay for styling purposes

const TemplateYN = `{{- Bold .Prompt -}} 
{{ if .Done }}
	{{- if .YesSelected -}}
		{{- print " " (Selected "Yes") -}}
	{{- else if .NoSelected -}}
		{{- print " " (Selected "No") -}}
	{{- end -}}

{{- else -}}
	{{- if .YesSelected -}}
		{{- print (Selected " [▸]Yes") " [ ]No" -}}
	{{- else if .NoSelected -}}
		{{- print " [ ]Yes " (Selected "[▸]No")  -}}
	{{- end -}}
{{ end }}
`

type model struct {
	confirm uimsg.Confirm
	value   bool
	done    bool

	promptTmpl *template.Template

	colorProfile   termenv.Profile
	selectionColor termenv.Color

	input  *io.Reader
	output *io.Writer

	width  int
	height int
	ready  bool

	help     help.Model
	keyMap   KeyMap
	viewport viewport.Model
}

func initialModel(c uimsg.Confirm) *model {
	return &model{
		confirm: c,
		value:   false,
		done:    false,
		help:    help.New(),
		keyMap:  defaultKeyMap(),
	}
}

func (m *model) Init() tea.Cmd {
	m.initTemplate()
	return nil
}

func (m *model) initTemplate() {
	p := termenv.ColorProfile()

	if tmpl, err := template.New("view").
		Funcs(termenv.TemplateFuncs(p)).
		Funcs(template.FuncMap{"Selected": func(values ...interface{}) string {
			s := termenv.String(values[0].(string))
			s = s.Foreground(m.selectionColor)
			s = s.Bold()
			return s.String()
		}}).
		Parse(TemplateYN); err != nil {
		panic(err)
	} else {
		m.promptTmpl = tmpl
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// cmd  tea.Cmd
	var cmds []tea.Cmd

	preValue := m.value

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			m.value = false
			m.done = true
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Apply):
			m.done = true
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Yes):
			m.value = true
		case key.Matches(msg, m.keyMap.No):
			m.value = false
		case key.Matches(msg, m.keyMap.Toggle):
			m.value = !m.value
		}

		if preValue != m.value {
			m.viewport.SetContent(m.content())
		}

	case tea.WindowSizeMsg:
		m.width = zeroAwareMin(msg.Width, 0)
		m.height = msg.Height
		m.ready = true

		height := msg.Height - lipgloss.Height(m.help.View(m)) - docStyle.GetMarginBottom() - docStyle.GetMarginTop()

		if !m.ready {
			m.viewport = viewport.Model{Width: msg.Width, Height: height}
			m.viewport.YPosition = 0
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = height
		}

		m.help.Width = msg.Width
		m.viewport.SetContent(m.content())
	}

	m.viewport, _ = m.viewport.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m model) content() string {
	var s string

	if !m.done {
		s += m.confirm.Header(m.width)
	}

	viewBuffer2 := &bytes.Buffer{}
	if err := m.promptTmpl.Execute(viewBuffer2, map[string]interface{}{
		"Prompt":      m.confirm.Prompt,
		"Done":        m.done,
		"YesSelected": m.value,
		"NoSelected":  !m.value,
	}); err != nil {
		panic(err)
	} else {
		s += viewBuffer2.String()
	}

	return m.wrap(s)
}

func (m *model) isScrollable() bool {
	return m.viewport.YOffset > 0 || m.viewport.ScrollPercent() < 1
}

func (m *model) View() string {
	var sections []string
	sections = append(sections, m.viewport.View())
	sections = append(sections, m.help.View(m))
	x := docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
	return x
}

func (m *model) wrap(text string) string {
	return wrap.String(wordwrap.String(text, m.width), m.width)
}

func Confirm(confirm uimsg.Confirm, options ...Option) bool {
	m := initialModel(confirm)
	for _, o := range options {
		o.apply(m)
	}

	var teaOptions []tea.ProgramOption
	if m.input != nil {
		teaOptions = append(teaOptions, tea.WithInput(*m.input))
	}
	if m.output != nil {
		teaOptions = append(teaOptions, tea.WithOutput(*m.output))
	}

	// teaOptions = append(teaOptions, tea.WithAltScreen())

	p := tea.NewProgram(m, teaOptions...)

	if err := p.Start(); err != nil {
		panic(err)
	}
	return m.value
}

type Option interface {
	apply(c *model)
}

type optionFunc func(o *model)

func (f optionFunc) apply(o *model) {
	f(o)
}

func WithIn(input io.Reader) Option {
	return optionFunc(func(c *model) {
		c.input = &input
	})
}

func WithOut(out io.Writer) Option {
	return optionFunc(func(c *model) {
		c.output = &out
	})
}

func WithSelectionColor(color string) Option {
	return optionFunc(func(c *model) {
		c.selectionColor = c.colorProfile.Color(color)
	})
}

func zeroAwareMin(a int, b int) int {
	switch {
	case a == 0:
		return b
	case b == 0:
		return a
	case a > b:
		return b
	default:
		return a
	}
}
