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

func newAppForm(parameters []model.Parameter) *appForm {
	return &appForm{
		app:            tview.NewApplication(),
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
			AddButton("Save", func() {
				a.success = true
				a.app.Stop()
			}).
			AddButton("Quit", func() { a.app.Stop() })
	case a.nextInputIndex < len(a.parameters):
		param := a.parameters[a.nextInputIndex]
		if len(param.Values) == 0 {
			a.form.AddInputField(padLength(param.Name, a.maxTitleLength), param.DefaultValue, defaultInputWidth, nil, nil)
		} else {
			a.form.AddFormItem(tview.NewInputField().
				SetLabel(padLength(param.Name, a.maxTitleLength)).
				SetText("value").
				SetFieldWidth(defaultInputWidth).
				SetText(param.DefaultValue).
				SetAutocompleteFunc(func(currentText string) (entries []string) {
					return param.Values
				}),
			)
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

func (c cliTerminal) ShowParameterForm(parameters []model.Parameter) ([]string, error) {
	return newAppForm(parameters).show()
}
