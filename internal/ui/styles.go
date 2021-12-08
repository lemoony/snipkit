package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var defaultStyle = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorReset,
	ContrastBackgroundColor:     tcell.ColorBlue,
	MoreContrastBackgroundColor: tcell.ColorGreen,
	BorderColor:                 tcell.ColorLightGrey,
	TitleColor:                  tcell.ColorYellow,
	GraphicsColor:               tcell.ColorGreenYellow,
	PrimaryTextColor:            tcell.ColorReset,
	SecondaryTextColor:          tcell.ColorYellow,
	TertiaryTextColor:           tcell.ColorGreen,
	InverseTextColor:            tcell.ColorBlue,
	ContrastSecondaryTextColor:  tcell.ColorDarkBlue,
}
