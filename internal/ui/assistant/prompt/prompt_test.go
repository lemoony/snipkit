package prompt

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

func TestShowPrompt_WithHistory(t *testing.T) {
	history := []string{"previous prompt 1", "previous prompt 2"}
	config := Config{History: history}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("[1] previous prompt 1")
		c.ExpectString("[2] previous prompt 2")
		c.Send("new prompt text")
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		ok, prompt := ShowPrompt(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.True(t, ok)
		assert.Equal(t, "new prompt text", prompt)
	})
}

func TestShowPrompt_EmptyInput(t *testing.T) {
	config := Config{History: []string{}}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("What do you want the script to do?")
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		ok, prompt := ShowPrompt(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.True(t, ok)
		assert.Equal(t, "", prompt)
	})
}

func TestShowPrompt_Cancel(t *testing.T) {
	config := Config{History: []string{"previous prompt"}}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("Do you want to provide additional context or change anything?")
		c.SendKey(termtest.KeyStrC)
	}, func(stdio termutil.Stdio) {
		ok, _ := ShowPrompt(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.False(t, ok)
	})
}
