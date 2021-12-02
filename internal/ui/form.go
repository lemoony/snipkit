package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lemoony/snippet-kit/internal/model"
)

var defaultStyle = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorDefault,
	ContrastBackgroundColor:     tcell.ColorDefault,
	MoreContrastBackgroundColor: tcell.ColorDefault,
	BorderColor:                 tcell.ColorDefault,
	TitleColor:                  tcell.ColorBlue, // X
	GraphicsColor:               tcell.ColorGray,
	PrimaryTextColor:            tcell.ColorRed,      // X
	SecondaryTextColor:          tcell.ColorGreen,    // shortcuts color
	TertiaryTextColor:           tcell.ColorDarkCyan, // subtitle // description color
	InverseTextColor:            tcell.ColorLightGrey,
	ContrastSecondaryTextColor:  tcell.ColorGrey,
}

type ParameterType int

type form struct {
	parameter             []model.Parameter
	currentParameterIndex int
	results               []string
}

func ShowParameterForm(parameters []model.Parameter) []string {
	tview.Styles = defaultStyle

	app := tview.NewApplication()
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	newForm := form{
		parameter:             parameters,
		currentParameterIndex: 0,
		results:               make([]string, len(parameters)),
	}

	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

	flex.AddItem(form, 0, 1, true)

	firstParam := newForm.parameter[newForm.currentParameterIndex]
	form.AddFormItem(tview.NewInputField().
		SetLabel(firstParam.Name).
		SetChangedFunc(func(text string) {
			focusedItem, _ := form.GetFocusedItemIndex()
			newForm.results[focusedItem] = text
		}).
		SetPlaceholder(firstParam.Description))

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			if focusedItem, _ := form.GetFocusedItemIndex(); focusedItem == len(newForm.parameter)-1 {
				app.Stop()
				return event
			} else if focusedItem == newForm.currentParameterIndex {
				newForm.currentParameterIndex += 1
				nextParam := newForm.parameter[newForm.currentParameterIndex]

				form.AddFormItem(tview.NewInputField().
					SetLabel(nextParam.Name).
					SetChangedFunc(func(text string) {
						focusedItem, _ := form.GetFocusedItemIndex()
						newForm.results[focusedItem] = text
					}).
					SetPlaceholder(nextParam.Description))
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}

	return newForm.results
}
