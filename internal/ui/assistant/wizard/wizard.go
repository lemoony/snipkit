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
	OptionTryAgain Option = iota
	OptionSaveExit
	OptionDontSaveExit

	formStepSelectOption formStep = iota
	formStepEnterFilename
	formStepEnterSnippetName
)

type Config struct {
	ProposedFilename    string
	ProposedSnippetName string
}

type Result struct {
	SelectedOption Option
	Filename       string
	SnippetTitle   string
}

type formModel struct {
	initialFilename    string
	initialSnippetName string

	quitting         bool
	success          bool
	option           Option
	filename         string
	snippetName      string
	step             formStep
	list             list.Model
	styler           style.Style
	filenameInput    textinput.Model
	snippetNameInput textinput.Model
}

func ShowAssistantWizard(config Config, styler style.Style, teaOptions ...tea.ProgramOption) (bool, Result) {
	m := newFormModel(config, styler)

	teaModel, err := tea.NewProgram(m, teaOptions...).Run()
	if err != nil {
		return false, Result{SelectedOption: OptionDontSaveExit}
	}

	model := teaModel.(*formModel)

	return model.success, Result{
		SelectedOption: model.option,
		Filename:       stringutil.StringOrDefault(model.filename, model.initialFilename),
		SnippetTitle:   stringutil.StringOrDefault(model.snippetName, model.initialSnippetName),
	}
}

func newFormModel(config Config, styler style.Style) *formModel {
	listModel := createListModel(styler)
	filenameInputModel := createFilenameInputModel(config, styler)
	snippetNameInputModel := createSnippetNameInputModel(config, styler)

	return &formModel{
		step:               formStepSelectOption,
		initialFilename:    config.ProposedFilename,
		initialSnippetName: config.ProposedSnippetName,
		filename:           config.ProposedFilename,
		snippetName:        config.ProposedSnippetName,
		list:               listModel,
		styler:             styler,
		filenameInput:      filenameInputModel,
		snippetNameInput:   snippetNameInputModel,
	}
}

func createListModel(styler style.Style) list.Model {
	items := []list.Item{
		listItem{title: "Tweak prompt and/or try again", option: OptionTryAgain},
		listItem{title: "Exit & Save", option: OptionSaveExit},
		listItem{title: "Exit & Don't save", option: OptionDontSaveExit},
	}

	listStyles := list.NewDefaultItemStyles()
	listStyles.SelectedTitle = lipgloss.NewStyle().
		SetString(">").
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.BorderColor().Value()).
		Foreground(styler.ActiveColor().Value()).
		PaddingLeft(1)

	listStyles.NormalTitle = lipgloss.NewStyle().
		SetString(" ").
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.BorderColor().Value()).
		PaddingLeft(1)

	const listWidth = 32
	const listHeight = 5

	l := list.New(items, list.DefaultDelegate{
		ShowDescription: false,
		Styles:          listStyles,
	}, listWidth, listHeight)

	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return l
}

func createFilenameInputModel(config Config, styler style.Style) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Type in filename..."
	input.SetValue(config.ProposedFilename)
	input.Focus()
	input.CharLimit = 256
	input.Width = 30
	input.PromptStyle = lipgloss.NewStyle().
		Foreground(styler.ActiveColor().Value()).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.BorderColor().Value()).
		PaddingLeft(1)
	input.Cursor.Style = lipgloss.NewStyle().Foreground(styler.HighlightColor().Value())

	return input
}

func createSnippetNameInputModel(config Config, styler style.Style) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Type in snippet name..."
	input.SetValue(config.ProposedSnippetName)
	input.CharLimit = 256
	input.Width = 30
	input.PromptStyle = lipgloss.NewStyle().
		Foreground(styler.ActiveColor().Value()).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.BorderColor().Value()).
		PaddingLeft(1)
	input.Cursor.Style = lipgloss.NewStyle().Foreground(styler.HighlightColor().Value())

	return input
}

func (m *formModel) Init() tea.Cmd {
	return m.list.StartSpinner()
}

func (m *formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case formStepSelectOption:
		return m.updateSelectOption(msg)
	case formStepEnterFilename:
		return m.updateEnterFilename(msg)
	case formStepEnterSnippetName:
		return m.updateEnterSnippetName(msg)
	}
	return m, nil
}

func (m *formModel) updateSelectOption(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			selectedItem := m.list.SelectedItem().(listItem)
			m.option = selectedItem.option
			if m.option == OptionSaveExit {
				m.step = formStepEnterFilename
				return m, m.filenameInput.Focus()
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
}

func (m *formModel) updateEnterFilename(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.filenameInput, cmd = m.filenameInput.Update(msg)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			m.filename = m.filenameInput.Value()
			m.step = formStepEnterSnippetName
			return m, m.snippetNameInput.Focus()
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m *formModel) updateEnterSnippetName(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.snippetNameInput, cmd = m.snippetNameInput.Update(msg)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			m.snippetName = m.snippetNameInput.Value()
			m.success = true
			m.quitting = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m *formModel) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(m.styler.Title("SnipKit Assistant"))
	sb.WriteString("\n")

	switch m.step {
	case formStepSelectOption:
		sb.WriteString(m.styler.PromptLabel("The snippet was executed. What now?"))
		sb.WriteString("\n")
		sb.WriteString(m.list.View())
	case formStepEnterFilename:
		sb.WriteString(m.styler.PromptLabel("Snippet Filename:"))
		sb.WriteString("\n")
		sb.WriteString(m.filenameInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(m.styler.InputHelp("Press enter to confirm"))
	case formStepEnterSnippetName:
		sb.WriteString(m.styler.PromptLabel("Snippet Filename: ") + m.filename)
		sb.WriteString("\n")
		sb.WriteString(m.styler.PromptLabel("Snippet Name:"))
		sb.WriteString("\n")
		sb.WriteString(m.snippetNameInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(m.styler.InputHelp("Press enter to confirm or ESC to quit"))
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
