package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lemoony/snippet-kit/internal/model"
)

const defaultInputWidth = 25

type OkButton string

const (
	OkButtonExecute = "Execute"
	OkButtonPrint   = "Print"
)

type appForm struct {
	app            *tview.Application
	form           *tview.Form
	parameters     []model.Parameter
	maxTitleLength int
	nextInputIndex int
	success        bool
	okButton       OkButton
}

func (c cliTerminal) ShowParameterForm(parameters []model.Parameter, okButton OkButton) ([]string, bool) {
	if len(parameters) == 0 {
		return []string{}, true
	}
	return newAppForm(parameters, c.screen, okButton).show()
}

func newAppForm(parameters []model.Parameter, screen tcell.Screen, okButton OkButton) *appForm {
	app := tview.NewApplication()
	app.SetScreen(screen)

	return &appForm{
		app:            app,
		form:           tview.NewForm(),
		parameters:     parameters,
		maxTitleLength: getMaxWidthParameter(parameters),
		nextInputIndex: 0,
		okButton:       okButton,
	}
}

func (a *appForm) show() ([]string, bool) {
	a.form.
		SetItemPadding(1).
		SetBorder(true).
		SetTitle("This snippet requires parameters").
		SetTitleAlign(tview.AlignLeft)

	a.addNextInput()

	a.form.SetLabelColor(currentTheme.parametersLabelColor())
	a.form.SetButtonBackgroundColor(currentTheme.selectedButtonBackgroundColor())
	a.form.SetFieldBackgroundColor(currentTheme.parametersFieldBackgroundColor())
	a.form.SetFieldTextColor(currentTheme.parametersFieldTextColor())

	if err := a.app.SetRoot(a.form, true).SetFocus(a.form).Run(); err != nil {
		panic(err)
	}

	if a.success {
		return a.collectParameterValues(), true
	}

	return []string{}, false
}

func (a *appForm) addNextInput() {
	switch {
	case a.nextInputIndex > len(a.parameters):
		return
	case a.nextInputIndex == len(a.parameters):
		a.form.
			AddButton(string(a.okButton), func() {
				a.success = true
				a.app.Stop()
			}).
			AddButton("Quit", func() {
				a.app.Stop()
			})

	case a.nextInputIndex < len(a.parameters):
		param := a.parameters[a.nextInputIndex]

		field := tview.NewInputField().
			SetLabel(padLength(param.Name, a.maxTitleLength)).
			SetText(param.DefaultValue).
			SetFieldWidth(defaultInputWidth).
			SetDoneFunc(a.fieldDoneFuc)

		if len(param.Values) > 0 {
			field.SetAutocompleteBackgroundColor(currentTheme.parametersAutocompleteBackgroundColor()).
				SetAutocompleteSelectBackgroundColor(currentTheme.parametersAutocompleteSelectedBackgroundColor()).
				SetAutocompleteMainTextColor(currentTheme.parametersAutocompleteTextColor()).
				SetAutocompleteSelectedTextColor(currentTheme.parametersAutocompleteSelectedTextColor()).
				SetAutocompleteFunc(func(currentText string) (entries []string) {
					return param.Values
				})
		}
		a.form.AddFormItem(field)
	}

	a.nextInputIndex++
}

func (a *appForm) fieldDoneFuc(key tcell.Key) {
	if key == tcell.KeyEnter || key == tcell.KeyTAB {
		a.addNextInput()
	}
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
