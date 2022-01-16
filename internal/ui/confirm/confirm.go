package confirm

import (
	"bytes"
	"io"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"github.com/muesli/termenv"
)

var fullScreenMargin = []int{1, 2, 0, 2}

type Key []string

var (
	keyYes    = Key{"y", "Y"}
	keyNo     = Key{"n", "N"}
	keyLeft   = Key{"left"}
	keyRight  = Key{"right"}
	keyToggle = Key{"tab"}
	keySubmit = Key{"enter", "ctrl+j"}
	keyAbort  = Key{"ctrl+c", "esc"}
)

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
	header  string
	confirm string
	value   bool
	done    bool

	fullscreen bool
	promptTmpl *template.Template

	colorProfile   termenv.Profile
	selectionColor termenv.Color

	input  *io.Reader
	output *io.Writer

	width int
}

func initialModel(prompt, header string) *model {
	return &model{
		header:  header,
		confirm: prompt,
		value:   false,
		done:    false,
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keyMatches(msg, keyAbort):
			m.value = false
			m.done = true
			return m, tea.Quit
		case keyMatches(msg, keySubmit):
			m.done = true
			return m, tea.Quit
		case keyMatches(msg, keyLeft):
			m.value = true
		case keyMatches(msg, keyRight):
			m.value = false
		case keyMatches(msg, keyToggle):
			m.value = !m.value
		case keyMatches(msg, keyYes):
			m.value = true
			m.done = true
			return m, tea.Quit
		case keyMatches(msg, keyNo):
			m.value = false
			m.done = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = zeroAwareMin(msg.Width, 0)
	}

	return m, nil
}

func (m model) View() string {
	var s string

	if !m.done {
		s += m.header + "\n"
	}

	viewBuffer2 := &bytes.Buffer{}
	if err := m.promptTmpl.Execute(viewBuffer2, map[string]interface{}{
		"Prompt":      m.confirm,
		"Done":        m.done,
		"YesSelected": m.value,
		"NoSelected":  !m.value,
	}); err != nil {
		panic(err)
	} else {
		s += viewBuffer2.String()
	}

	content := m.wrap(s)

	if m.fullscreen {
		content = lipgloss.NewStyle().Margin(fullScreenMargin...).SetString(content).String()
	}

	return content
}

func (m *model) wrap(text string) string {
	return wrap.String(wordwrap.String(text, m.width), m.width)
}

func keyMatches(key tea.KeyMsg, expected Key) bool {
	for _, m := range expected {
		if m == key.String() {
			return true
		}
	}
	return false
}

func Confirm(prompt string, header string, options ...Option) bool {
	m := initialModel(prompt, header)
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
	if m.fullscreen {
		teaOptions = append(teaOptions, tea.WithAltScreen())
	}

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

func WithFullscreen() Option {
	return optionFunc(func(c *model) {
		c.fullscreen = true
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
