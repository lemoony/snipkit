package wizard

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

type (
	Option   int
	formStep int
)

const (
	OptionNone         = Option(0)
	OptionTryAgain     = Option(1)
	OptionSaveExit     = Option(2)
	OptionDontSaveExit = Option(3)

	formStepSelectOption  = formStep(0)
	formStepEnterFilename = formStep(1)
)

type Config struct {
	ProposedFilename string
}

type Result struct {
	SelectedOption Option
	Filename       string
}

type formModel struct {
	initialFilename string

	quitting bool
	success  bool
	option   Option
	filename string
	step     formStep
	list     list.Model
	styler   style.Style
	input    textinput.Model
}

func ShowAssistantWizard(config Config, styler style.Style) (bool, Result) {
	m := newModel(config, styler)

	teaModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return false, Result{SelectedOption: OptionDontSaveExit}
	}

	model := teaModel.(*formModel)

	return model.success, Result{
		SelectedOption: model.option,
		Filename:       stringutil.StringOrDefault(model.filename, model.initialFilename),
	}
}

func newModel(config Config, styler style.Style) *formModel {
	items := []list.Item{
		listItem{title: "Tweak prompt and/or try again", option: OptionTryAgain},
		listItem{title: "Exit & Save", option: OptionSaveExit},
		listItem{title: "Exit & Don't save", option: OptionDontSaveExit},
	}

	styls := list.NewDefaultItemStyles()
	styls.SelectedTitle = lipgloss.NewStyle().
		SetString(">").
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.BorderColor().Value()).
		Foreground(styler.ActiveColor().Value()).
		PaddingLeft(1)

	styls.NormalTitle = lipgloss.NewStyle().SetString(" ").
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.BorderColor().Value()).
		PaddingLeft(1)

	l := list.New(items, list.DefaultDelegate{
		ShowDescription: false,
		Styles:          styls,
	}, 0, 15)

	l.SetSize(32, 5)

	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	input := textinput.New()
	input.Placeholder = "Type in filename..."
	input.SetValue(config.ProposedFilename)
	input.Focus()
	input.CharLimit = 256
	input.Width = 30
	input.PromptStyle = lipgloss.NewStyle().Foreground(styler.ActiveColor().Value())
	input.Cursor.Style = lipgloss.NewStyle().Foreground(styler.HighlightColor().Value())
	input.PromptStyle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.BorderColor().Value()).
		Foreground(styler.ActiveColor().Value()).
		PaddingLeft(1)

	return &formModel{
		option:          OptionNone,
		step:            formStepSelectOption,
		filename:        config.ProposedFilename,
		initialFilename: config.ProposedFilename,
		list:            l,
		styler:          styler,
		input:           input,
	}
}

func (m *formModel) Init() tea.Cmd {
	return m.list.StartSpinner()
}

func (m *formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case formStepSelectOption:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				selectedItem := m.list.SelectedItem().(listItem)
				m.option = selectedItem.option
				if m.option == OptionSaveExit {
					m.step = formStepEnterFilename
					return m, m.input.Focus()
				}
				m.success = true
				m.quitting = true
				return m, tea.Quit
			case tea.KeyCtrlC, tea.KeyEsc:
				m.quitting = true
				return m, tea.Quit
			}
		}
		return m, cmd

	case formStepEnterFilename:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				m.filename = m.input.Value()
				m.quitting = true
				m.success = true
				return m, tea.Quit
			case tea.KeyCtrlC, tea.KeyEsc:
				m.quitting = true
				return m, tea.Quit
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m *formModel) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(m.styler.Title("SnipKit Assistant"))
	sb.WriteString("\n")

	switch m.step {
	case formStepSelectOption:
		sb.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(m.styler.TitleColor().Value()).
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderForeground(m.styler.BorderColor().Value()).
			PaddingLeft(1).
			Render("The snippet was executed. What now?"))

		sb.WriteString("\n")
		sb.WriteString(m.list.View())
	case formStepEnterFilename:

		sb.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(m.styler.TitleColor().Value()).
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderForeground(m.styler.BorderColor().Value()).
			PaddingLeft(1).
			Render("Snippet Filename:"))

		sb.WriteString("\n")
		sb.WriteString(m.input.View())
		sb.WriteString("\n\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(m.styler.PlaceholderColor().Value()).Render("Press enter to confirm"))
	}

	return sb.String()
}

type listItem struct {
	title  string
	option Option
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return i.title }
