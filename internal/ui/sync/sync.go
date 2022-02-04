package sync

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
)

const (
	managerMarginLeft = 3
)

var (
	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).
			PaddingRight(1).
			String()

	continueKeyBinding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "apply"),
	)

	quitKeyBinding = key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)
)

type Screen struct {
	program *tea.Program
}

type UpdateStateMsg struct {
	Status       appModel.SyncStatus
	ManagerState *ManagerState
}

type state struct {
	Status       appModel.SyncStatus
	ManagerState map[appModel.ManagerKey]ManagerState
}

func (s *state) isWaitingForInput() (chan appModel.SyncInputResult, bool) {
	for _, m := range s.ManagerState {
		if m.Input != nil && m.Input.Input != nil {
			return m.Input.Input, true
		}
	}
	return nil, false
}

type ManagerState struct {
	Key    appModel.ManagerKey
	Status appModel.SyncStatus
	Lines  []appModel.SyncLine
	Input  *appModel.SyncInput
	Error  error
}

type model struct {
	styler style.Style

	input  *io.Reader
	output *io.Writer

	state state

	width  int
	height int

	spinner   spinner.Model
	textinput textinput.Model
}

func New(options ...Option) *Screen {
	m := &model{state: state{
		Status:       appModel.SyncStatusFinished,
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

	return &Screen{
		program: tea.NewProgram(m, teaOptions...),
	}
}

func (s *Screen) Start() {
	if err := s.program.Start(); err != nil {
		panic(err)
	}
}

func (s *Screen) Send(msg UpdateStateMsg) {
	s.program.Send(msg)
}

func (m *model) Init() tea.Cmd {
	return m.spinner.Tick
}

//nolint:gocognit,gocyclo //TODO refactor at a later point
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case UpdateStateMsg:
		if status := msg.Status; status != 0 {
			m.state.Status = status
			if status == appModel.SyncStatusFinished || status == appModel.SyncStatusAborted {
				return m, tea.Quit
			}
		}
		if managerState := msg.ManagerState; managerState != nil {
			m.state.ManagerState[managerState.Key] = *managerState

			if login := managerState.Input; login != nil && login.Type == appModel.SyncLoginTypeText {
				m.textinput = textinput.New()
				m.textinput.Placeholder = login.Placeholder
				m.textinput.Focus()
			}

			if managerState.Status == appModel.SyncStatusAborted {
				return m, tea.Quit
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, quitKeyBinding):
			if c, isWaiting := m.state.isWaitingForInput(); isWaiting {
				c <- appModel.SyncInputResult{Abort: true}
			}
		case key.Matches(msg, continueKeyBinding):
			if c, isWaiting := m.state.isWaitingForInput(); isWaiting {
				if m.textinput.Focused() {
					c <- appModel.SyncInputResult{Text: m.textinput.Value()}
				} else {
					c <- appModel.SyncInputResult{Continue: true}
				}
			}
		default:
			if m.textinput, cmd = m.textinput.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	sections := []string{"Syncing all managers..."}

	for _, v := range m.state.ManagerState {
		sections = append(sections, fmt.Sprintf("%s Syncing %s...", m.spinner.View(), string(v.Key)))

		for _, l := range v.Lines {
			sections = append(sections, lipgloss.NewStyle().MarginLeft(managerMarginLeft).Render(l.Value))
		}

		if input := v.Input; input != nil {
			sections = append(
				sections,
				lipgloss.NewStyle().MarginLeft(managerMarginLeft).MarginTop(1).Render(input.Content),
			)

			if input.Type == appModel.SyncLoginTypeText {
				sections = append(sections, lipgloss.NewStyle().MarginLeft(managerMarginLeft).Render(m.textinput.View()))
			}
		}
	}

	if m.state.Status == appModel.SyncStatusFinished {
		sections = append(sections, fmt.Sprintf("%s All done.\n", checkMark))
	}

	return m.wrap(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m *model) wrap(text string) string {
	y := wordwrap.String(text, m.width)
	return wrap.String(y, m.width)
}
