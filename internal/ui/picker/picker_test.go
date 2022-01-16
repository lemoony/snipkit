package picker

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

func Test_ShowPicker(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Which snippet manager should be added to your configuration")
		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyUp)
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		index, ok := ShowPicker([]Item{
			NewItem("title1", "desc1"),
			NewItem("title2", "desc2"),
			NewItem("title3", "desc3"),
		}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.Equal(t, 1, index)
		assert.True(t, ok)
	})
}
