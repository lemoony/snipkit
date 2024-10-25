package spinner

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

func Test_ShowSpinner(t *testing.T) {
	tests := []struct {
		name       string
		exitByChan bool
		stopKey    termtest.Key
	}{
		{name: "stop by chan", stopKey: termtest.KeyTab, exitByChan: true},
		{name: "stop by esc", stopKey: termtest.KeyEsc, exitByChan: false},
		{name: "stop by str+c", stopKey: termtest.KeyStrC, exitByChan: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopChan := make(chan bool)
			termtest.RunTerminalTest(t, func(c *termtest.Console) {
				time.Sleep(100 * time.Millisecond)
				c.ExpectString("Test text...")
				time.Sleep(50 * time.Millisecond)
				c.SendKey(tt.stopKey)
				if tt.exitByChan {
					stopChan <- true
				}
			}, func(stdio termutil.Stdio) {
				ShowSpinner("Test text...", stopChan, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
			})
		})
	}
}
