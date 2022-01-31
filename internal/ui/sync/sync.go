package sync

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
)

var checkMark = lipgloss.NewStyle().SetString("✓").
	Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).
	PaddingRight(1).
	String()

type State struct {
	Done         bool
	ManagerState map[appModel.ManagerKey]ManagerState
}

type ManagerState struct {
	Key        appModel.ManagerKey
	InProgress bool
	Error      error
}

type UpdateStateMsg struct {
	State State
}

type model struct {
	styler style.Style

	input  *io.Reader
	output *io.Writer

	state State

	width  int
	height int

	spinner spinner.Model
}

func Show(options ...Option) *tea.Program {
	m := &model{state: State{
		Done:         false,
		ManagerState: map[appModel.ManagerKey]ManagerState{},
	}}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m.spinner = s

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

	return tea.NewProgram(m, teaOptions...)
}

func (m *model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case UpdateStateMsg:
		if msg.State.Done {
			m.state.Done = true
			cmd = tea.Quit
		}
	case ManagerState:
		m.state.ManagerState[msg.Key] = msg
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m *model) View() string {
	sections := []string{"Syncing all managers..."}

	for _, v := range m.state.ManagerState {
		if v.InProgress {
			sections = append(sections, fmt.Sprintf("%s Syncing %s...", m.spinner.View(), string(v.Key)))
		} else {
			sections = append(sections, fmt.Sprintf("%s Syncing %s... done", checkMark, string(v.Key)))
		}
	}

	if m.state.Done {
		sections = append(sections, fmt.Sprintf("%s All done.\n", checkMark))
	}

	return m.wrap(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m *model) wrap(text string) string {
	return wrap.String(wordwrap.String(text, m.width), m.width)
}
