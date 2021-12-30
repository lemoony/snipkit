package ui

import (
	"io/ioutil"
	"os"
	"path"
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

func Test_Confirm(t *testing.T) {
	runTest(t, func(c *expect.Console) {
		_, _ = c.ExpectString("Are you sure? (y/N)")
		_, _ = c.SendLine("Y")
		_, _ = c.ExpectEOF()
	}, func(stdio terminal.Stdio) {
		term := ActualCLI{stdio: stdio}
		confirmed, err := term.Confirm("Are you sure?")
		assert.NoError(t, err)
		assert.True(t, confirmed)
	})
}

func Test_OpenEditor(t *testing.T) {
	runTest(t, func(c *expect.Console) {
		_, _ = c.Send("iHello world\x1b")
		time.Sleep(time.Second)
		_, _ = c.SendLine(":wq!")
	}, func(stdio terminal.Stdio) {
		term := ActualCLI{stdio: stdio}

		testFile := path.Join(t.TempDir(), "testfile")
		_, err := os.Create(testFile)
		assert.NoError(t, err)

		err = term.OpenEditor(testFile, "")
		assert.NoError(t, err)
		bytes, err := ioutil.ReadFile(testFile) //nolint:gosec // potential file inclusion via variable
		assert.NoError(t, err)
		assert.Equal(t, "Hello world\n", string(bytes))
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
		{name: "default editor windows", disabled: runtime.GOOS != windows, expected: defaultEditor},
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
					assert.NoError(t, os.Unsetenv(tt.envEditor))
				} else {
					assert.NoError(t, os.Setenv("EDITOR", tt.envEditor))
				}
				if tt.envVisual == "" {
					assert.NoError(t, os.Unsetenv(tt.envVisual))
				} else {
					assert.NoError(t, os.Setenv("VISUAL", tt.envVisual))
				}

				defer func() {
					_ = os.Unsetenv("VISUAL")
					_ = os.Unsetenv("EDITOR")
				}()

				editor := getEditor(tt.preferred)
				assert.Equal(t, tt.expected, editor)
			}
		})
	}
}
