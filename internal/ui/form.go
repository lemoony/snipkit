package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lemoony/snippet-kit/internal/model"
)

const defaultInputWidth = 25

type appForm struct {
	app            *tview.Application
	form           *tview.Form
	parameters     []model.Parameter
	maxTitleLength int
	nextInputIndex int
	success        bool
}

func (c cliTerminal) ShowParameterForm(parameters []model.Parameter) []string {
	if parameters, err := newAppForm(parameters, c.screen).show(); err != nil {
		panic(err)
	} else {
		return parameters
	}
}

func newAppForm(parameters []model.Parameter, screen tcell.Screen) *appForm {
	app := tview.NewApplication()
	app.SetScreen(screen)

	return &appForm{
		app:            app,
		form:           tview.NewForm(),
		parameters:     parameters,
		maxTitleLength: getMaxWidthParameter(parameters),
		nextInputIndex: 0,
	}
}

func (a *appForm) show() ([]string, error) {
	a.form.
		SetItemPadding(1).
		SetBorder(true).
		SetTitle("This snippet requires parameters").
		SetTitleAlign(tview.AlignLeft).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEnter, tcell.KeyTAB:
				a.addNextInput()
			}
			return event
		})

	a.addNextInput()

	a.form.SetLabelColor(currentTheme.parametersLabelColor())
	a.form.SetButtonBackgroundColor(currentTheme.selectedButtonBackgroundColor())
	a.form.SetFieldBackgroundColor(currentTheme.parametersFieldBackgroundColor())
	a.form.SetFieldTextColor(currentTheme.parametersFieldTextColor())

	if err := a.app.SetRoot(a.form, true).SetFocus(a.form).Run(); err != nil {
		return nil, err
	}

	if a.success {
		return a.collectParameterValues(), nil
	}

	return []string{}, nil
}

func (a *appForm) addNextInput() {
	switch {
	case a.nextInputIndex > len(a.parameters):
		return
	case a.nextInputIndex == len(a.parameters):
		a.form.
			AddButton("Execute", func() {
				a.success = true
				a.app.Stop()
			}).
			AddButton("Quit", func() {
				a.app.Stop()
			})

	case a.nextInputIndex < len(a.parameters):
		param := a.parameters[a.nextInputIndex]
		if len(param.Values) == 0 {
			a.form.AddInputField(padLength(param.Name, a.maxTitleLength), param.DefaultValue, defaultInputWidth, nil, nil)
		} else {
			field := tview.NewInputField().
				SetLabel(padLength(param.Name, a.maxTitleLength)).
				SetText("value").
				SetFieldWidth(defaultInputWidth).
				SetText(param.DefaultValue).
				SetAutocompleteBackgroundColor(currentTheme.parametersAutocompleteBackgroundColor()).
				SetAutocompleteSelectBackgroundColor(currentTheme.parametersAutocompleteSelectedBackgroundColor()).
				SetAutocompleteMainTextColor(currentTheme.parametersAutocompleteTextColor()).
				SetAutocompleteSelectedTextColor(currentTheme.parametersAutocompleteSelectedTextColor()).
				SetAutocompleteFunc(func(currentText string) (entries []string) {
					return param.Values
				})

			a.form.AddFormItem(field)
		}
	}

	a.nextInputIndex++
}

func (a *appForm) collectParameterValues() []string {
	results := make([]string, len(a.parameters))
	for i := range a.parameters {
		results[i] = a.form.GetFormItem(i).(*tview.InputField).GetText()
	}
	return results
}

func getMaxWidthParameter(parameters []model.Parameter) int {
	max := 0
	for i := range parameters {
		if l := len(parameters[i].Name); l > max {
			max = l
		}
	}
	return max
}

func padLength(title string, targetLength int) string {
	for l := len(title); l <= targetLength; l++ {
		title += " "
	}
	return title
}
