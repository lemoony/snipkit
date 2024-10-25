package prompt

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

func TestShowPrompt(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("placeholder text")
		c.Send("some input text")
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		ok, prompt := ShowPrompt("placeholder text", tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.True(t, ok)
		assert.Equal(t, "some input text", prompt)
	})
}

func TestShowPrompt_Cancel(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("placeholder text")
		c.SendKey(termtest.KeyStrC)
	}, func(stdio termutil.Stdio) {
		ok, _ := ShowPrompt("placeholder text", tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.False(t, ok)
	})
}
