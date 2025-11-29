package picker

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
)

func Test_ShowPicker(t *testing.T) {
	items := []Item{
		NewItem("title1", "desc1"),
		NewItem("title2", "desc2"),
		NewItem("title3", "desc3"),
	}

	m := NewModel("Which snippet manager should be added to your configuration", items, nil, style.NoopStyle)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for the title to appear
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Which snippet manager"))
	}, teatest.WithDuration(2*time.Second))

	// Navigate: down, down, up (should end up on item 1)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyUp})

	// Press enter to select
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Wait for the program to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	// Get the final model and check the selection
	finalModel := tm.FinalModel(t)
	pickerModel, ok := finalModel.(*Model)
	assert.True(t, ok, "expected *Model type")

	index, selected := pickerModel.SelectedIndex()
	assert.True(t, selected)
	assert.Equal(t, 1, index)
}

func Test_ShowPicker_Cancel(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyType
	}{
		{name: "esc", key: tea.KeyEsc},
		{name: "ctrl+c", key: tea.KeyCtrlC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := []Item{NewItem("title1", "desc1")}
			m := NewModel("Which snippet manager should be added to your configuration", items, nil, style.NoopStyle)
			tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

			// Wait for the title to appear
			teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
				return bytes.Contains(bts, []byte("Which snippet manager"))
			}, teatest.WithDuration(2*time.Second))

			// Send cancel key
			tm.Send(tea.KeyMsg{Type: tt.key})

			// Wait for the program to finish
			tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

			// Get the final model and check no selection was made
			finalModel := tm.FinalModel(t)
			pickerModel, ok := finalModel.(*Model)
			assert.True(t, ok, "expected *Model type")

			index, selected := pickerModel.SelectedIndex()
			assert.False(t, selected)
			assert.Equal(t, -1, index)
		})
	}
}

func Test_ShowPicker_WithPreselectedItem(t *testing.T) {
	items := []Item{
		NewItem("title1", "desc1"),
		NewItem("title2", "desc2"),
		NewItem("title3", "desc3"),
	}

	// Pre-select the second item
	preselected := items[1]
	m := NewModel("Select an item", items, &preselected, style.NoopStyle)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for the title to appear
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Select an item"))
	}, teatest.WithDuration(2*time.Second))

	// Press enter immediately (should select the pre-selected item)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Wait for the program to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	// Get the final model and check the selection
	finalModel := tm.FinalModel(t)
	pickerModel, ok := finalModel.(*Model)
	assert.True(t, ok, "expected *Model type")

	index, selected := pickerModel.SelectedIndex()
	assert.True(t, selected)
	assert.Equal(t, 1, index) // Should be index 1 (the pre-selected item)
}
