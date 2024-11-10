package wizard

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

type Option int

type formStep int

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
	selectionForm   *huh.Form
	filenameForm    *huh.Form
	initialFilename string

	quitting bool
	option   Option
	filename string
	step     formStep
}

func ShowAssistantWizard(config Config) Result {
	// TODO write test

	m := newModel(config)

	teaModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return Result{SelectedOption: OptionDontSaveExit}
	}

	model := teaModel.(*formModel)
	return Result{
		SelectedOption: model.option,
		Filename:       stringutil.StringOrDefault(model.filename, model.initialFilename),
	}
}

func newModel(config Config) *formModel {
	model := &formModel{
		option:          OptionNone,
		step:            formStepSelectOption,
		filename:        config.ProposedFilename,
		initialFilename: config.ProposedFilename,
	}

	// Create selection form
	s := huh.NewSelect[Option]().
		Title("The snippet was executed. What now?").
		Options(
			huh.NewOption("Try again and/or tweak prompt", OptionTryAgain),
			huh.NewOption("Exit & Save", OptionSaveExit),
			huh.NewOption("Exit & Don't save", OptionDontSaveExit),
		).
		Value(&model.option)

	model.selectionForm = huh.NewForm(huh.NewGroup(s))

	// Create filename form
	f := huh.NewInput().
		Title("Snippet filename:").
		Placeholder("Type in filename...").
		Value(&model.filename)

	model.filenameForm = huh.NewForm(huh.NewGroup(f))

	return model
}

func (m *formModel) Init() tea.Cmd {
	return m.selectionForm.Init()
}

func (m *formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case formStepSelectOption:
		form, cmd := m.selectionForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.selectionForm = f
		}

		if m.selectionForm.State == huh.StateCompleted {
			if m.option == OptionSaveExit {
				// Move to filename step
				m.step = formStepEnterFilename
				return m, m.filenameForm.Init()
			}
			m.quitting = true
			return m, tea.Quit
		}

		return m, cmd

	case formStepEnterFilename:
		form, cmd := m.filenameForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.filenameForm = f
		}

		if m.filenameForm.State == huh.StateCompleted {
			m.quitting = true
			return m, tea.Quit
		}

		return m, cmd
	}

	return m, nil
}

func (m *formModel) View() string {
	if m.quitting {
		return ""
	}

	switch m.step {
	case formStepSelectOption:
		return fmt.Sprintf("\n%s", m.selectionForm.View())
	case formStepEnterFilename:
		return fmt.Sprintf("\n%s", m.filenameForm.View())
	}

	return ""
}
