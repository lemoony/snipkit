package ui

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func Test_PrintMessage(t *testing.T) {
	runTest(t, func(c *expect.Console) {
		_, err := c.ExpectString("Hello world")
		assert.NoError(t, err)
		_, err = c.ExpectEOF()
		assert.NoError(t, err)
	}, func(stdio terminal.Stdio) {
		term := NewTerminal(WithStdio(stdio))
		term.PrintMessage("Hello world")
		time.Sleep(time.Millisecond * 100)
	})
}

func Test_PrintError(t *testing.T) {
	runTest(t, func(c *expect.Console) {
		_, err := c.ExpectString("Some error message")
		assert.NoError(t, err)
		_, err = c.ExpectEOF()
		assert.NoError(t, err)
	}, func(stdio terminal.Stdio) {
		term := NewTerminal(WithStdio(stdio))
		term.PrintError("Some error message")
		time.Sleep(time.Millisecond * 100)
	})
}

func Test_Confirm(t *testing.T) {
	runTest(t, func(c *expect.Console) {
		_, err := c.ExpectString("Are you sure? (y/N)")
		assert.NoError(t, err)
		_, err = c.SendLine("Y")
		assert.NoError(t, err)
		_, err = c.ExpectEOF()
		assert.NoError(t, err)
	}, func(stdio terminal.Stdio) {
		term := NewTerminal(WithStdio(stdio))
		confirmed, err := term.Confirm("Are you sure?")
		assert.NoError(t, err)
		assert.True(t, confirmed)
	})
}

func Test_getEditor(t *testing.T) {
	tests := []struct {
		name      string
		disabled  bool
		envVisual string
		envEditor string
		preferred string
		expected  string
	}{
		{name: "default editor unix", disabled: runtime.GOOS == windows, expected: defaultEditor},
		{name: "default editor windows", disabled: runtime.GOOS != windows, expected: defaultEditorWindows},
		{name: "editor env set", envEditor: "foo-editor", expected: "foo-editor"},
		{name: "visual env set", envVisual: "some-editor", expected: "some-editor"},
		{name: "editor + visual env set", envEditor: "some-editor", envVisual: "foo-editor", expected: "foo-editor"},
		{
			name:      "preferred editor set",
			envEditor: "foo-editor",
			envVisual: "foo-editor",
			preferred: "another-editor",
			expected:  "another-editor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.disabled {
				t.Skipf("Test %s is not enabled for this platform", tt.name)
			} else {
				if tt.envEditor == "" {
					_ = os.Unsetenv(envEditor)
				} else {
					_ = os.Setenv(envEditor, tt.envEditor)
				}
				if tt.envVisual == "" {
					_ = os.Unsetenv(envVisual)
				} else {
					_ = os.Setenv(envVisual, tt.envVisual)
				}

				defer func() {
					_ = os.Unsetenv(envVisual)
					_ = os.Unsetenv(envEditor)
				}()

				editor := getEditor(tt.preferred)
				assert.Equal(t, tt.expected, editor)
			}
		})
	}
}
