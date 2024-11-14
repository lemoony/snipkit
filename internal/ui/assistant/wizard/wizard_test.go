package wizard

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

func TestShowAssistantWizard_TryAgainOption(t *testing.T) {
	config := Config{ProposedFilename: "initial-filename.txt"}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("The snippet was executed. What now?")
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		success, result := ShowAssistantWizard(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.True(t, success)
		assert.Equal(t, OptionTryAgain, result.SelectedOption)
	})
}

func TestShowAssistantWizard_SaveExit_Default(t *testing.T) {
	config := Config{ShowSaveOption: true, ProposedFilename: "initial-filename.txt", ProposedSnippetName: "Foo title"}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("The snippet was executed. What now?")
		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyEnter)
		c.ExpectString("Snippet Filename:")
		c.SendKey(termtest.KeyEnter)
		c.ExpectString("Snippet Name:")
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		success, result := ShowAssistantWizard(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.True(t, success)
		assert.Equal(t, OptionSaveExit, result.SelectedOption)
		assert.Equal(t, config.ProposedFilename, result.Filename)
		assert.Equal(t, config.ProposedSnippetName, result.SnippetTitle)
	})
}

func TestShowAssistantWizard_SaveExit_Edit(t *testing.T) {
	config := Config{ShowSaveOption: true, ProposedFilename: "x.txt", ProposedSnippetName: "foo title"}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("The snippet was executed. What now?")
		c.SendKey(termtest.KeyDown) // Move to "Exit & Save"
		c.SendKey(termtest.KeyEnter)
		c.ExpectString("Snippet Filename:")
		for range config.ProposedFilename {
			c.SendKey(termtest.KeyDelete)
		}
		c.Send("new-filename.txt")
		c.SendKey(termtest.KeyEnter)
		c.ExpectString("Snippet Name:")
		for range config.ProposedSnippetName {
			c.SendKey(termtest.KeyDelete)
		}
		c.Send("example snippet")
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		success, result := ShowAssistantWizard(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.True(t, success)
		assert.Equal(t, OptionSaveExit, result.SelectedOption)
		assert.Equal(t, "new-filename.txt", result.Filename)
		assert.Equal(t, "example snippet", result.SnippetTitle)
	})
}

func TestShowAssistantWizard_DontSaveExit(t *testing.T) {
	config := Config{ShowSaveOption: true, ProposedFilename: "initial-filename.txt"}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("The snippet was executed. What now?")
		c.SendKey(termtest.KeyDown) // Move to "Exit & Save"
		c.SendKey(termtest.KeyDown) // Move to "Exit & Don't save"
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		success, result := ShowAssistantWizard(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.True(t, success)
		assert.Equal(t, OptionDontSaveExit, result.SelectedOption)
	})
}

func TestShowAssistantWizard_Cancel(t *testing.T) {
	config := Config{ProposedFilename: "initial-filename.txt"}

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("SnipKit Assistant")
		c.ExpectString("The snippet was executed. What now?")
		c.SendKey(termtest.KeyStrC)
	}, func(stdio termutil.Stdio) {
		success, _ := ShowAssistantWizard(config, style.Style{}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.False(t, success)
	})
}
