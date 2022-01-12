package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
)

func Test_Form_MultipleParameters(t *testing.T) {
	parameters := []model.Parameter{
		{Key: "VAR1", Name: "First", Description: "First parameter"},
		{Key: "VAR2", Name: "Second", Description: "Second parameter"},
	}

	runScreenTest(t, func(s tcell.Screen) {
		term := NewTerminal(WithScreen(s))
		values, ok := term.ShowParameterForm(parameters, OkButtonExecute)
		assert.True(t, ok)
		assert.Len(t, values, 2)
		assert.Equal(t, "First", values[0])
		assert.Equal(t, "Second", values[1])
	}, func(screen tcell.SimulationScreen) {
		sendString(t, "First", screen)
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))

		sendString(t, "Second", screen)
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))

		// hit ok button
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))
	})
}

func Test_Form_EnumParameter(t *testing.T) {
	parameters := []model.Parameter{
		{Key: "VAR1", Name: "First", Description: "First parameter", Values: []string{"FIRST_VAL", "SECOND_VAL"}},
	}

	runScreenTest(t, func(s tcell.Screen) {
		term := NewTerminal(WithScreen(s))
		values, ok := term.ShowParameterForm(parameters, OkButtonPrint)
		assert.True(t, ok)
		assert.Len(t, values, 1)
		assert.Equal(t, "SECOND_VAL", values[0])
	}, func(screen tcell.SimulationScreen) {
		sendString(t, "SECO", screen)
		// hint enter to select SECOND_VAL from list
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))
		// hint enter to exit input field
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))

		// hit ok button
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))
	})
}

func Test_Form_DefaultValue(t *testing.T) {
	parameters := []model.Parameter{
		{Key: "VAR1", Name: "First", Description: "First parameter", DefaultValue: "default-value"},
	}

	runScreenTest(t, func(s tcell.Screen) {
		term := NewTerminal(WithScreen(s))
		values, ok := term.ShowParameterForm(parameters, OkButtonPrint)
		assert.True(t, ok)
		assert.Len(t, values, 1)
		assert.Equal(t, "default-value", values[0])
	}, func(screen tcell.SimulationScreen) {
		// hint enter to exit input field
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))

		// hit ok button
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))
	})
}

func Test_Form_NoParameters(t *testing.T) {
	var parameters []model.Parameter

	term := NewTerminal()
	values, ok := term.ShowParameterForm(parameters, OkButtonPrint)
	assert.Len(t, values, 0)
	assert.True(t, ok)
}
