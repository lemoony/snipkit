package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var defaultThemeValues = ThemeValues{
	Background:             "",
	ContrastBackground:     "lightGrey",
	MoreContrastBackground: "green",
	Border:                 "lightGrey",
	Title:                  "yellow",
	Graphics:               "greenYellow",
	Text:                   "",
	SecondaryText:          "yellow",
	TertiaryText:           "green",
	InverseText:            "blue",
	ContrastSecondaryText:  "darkBlue",
}

func SetTheme(theme ThemeValues) {
	tview.Styles = toTviewTheme(theme)
}

func toTviewTheme(values ThemeValues) tview.Theme {
	return tview.Theme{
		PrimitiveBackgroundColor:    tcell.GetColor(values.Background),
		ContrastBackgroundColor:     tcell.GetColor(values.ContrastBackground),
		MoreContrastBackgroundColor: tcell.GetColor(values.MoreContrastBackground),
		BorderColor:                 tcell.GetColor(values.Border),
		TitleColor:                  tcell.GetColor(values.Title),
		GraphicsColor:               tcell.GetColor(values.Graphics),
		PrimaryTextColor:            tcell.GetColor(values.Text),
		SecondaryTextColor:          tcell.GetColor(values.SecondaryText),
		TertiaryTextColor:           tcell.GetColor(values.TertiaryText),
		InverseTextColor:            tcell.GetColor(values.InverseText),
		ContrastSecondaryTextColor:  tcell.GetColor(values.ContrastBackground),
	}
}
