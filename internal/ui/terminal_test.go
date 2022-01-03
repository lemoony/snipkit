package ui

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/model"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func Test_PrintMessage(t *testing.T) {
	runExpectTest(t, func(c *expect.Console) {
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
	runExpectTest(t, func(c *expect.Console) {
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
	runExpectTest(t, func(c *expect.Console) {
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

func Test_ShowLookup(t *testing.T) {
	snippets := []model.Snippet{
		{
			Title:    "Title 1",
			Content:  "Content: One",
			Language: model.LanguageYAML,
		},
		{
			Title:    "Title 2",
			Content:  "Content: Two",
			Language: model.LanguageYAML,
		},
	}

	runScreenTest(t, func(s tcell.Screen) {
		term := NewTerminal(WithScreen(s))
		selected, err := term.ShowLookup(snippets)

		assert.NoError(t, err)
		assert.Equal(t, 1, selected)
	}, func(screen tcell.SimulationScreen) {
		time.Sleep(time.Millisecond * 50)
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)))

		time.Sleep(time.Millisecond * 50)
		previewContent := getPreviewContents(screen)
		assert.Equal(t, snippets[1].Content, previewContent)

		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))
	})
}

func runScreenTest(t *testing.T, procedure func(s tcell.Screen), test func(s tcell.SimulationScreen)) {
	t.Helper()
	screen := mkTestScreen(t)

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		time.Sleep(time.Millisecond * 50)
		test(screen)
	}()

	procedure(screen)
	<-donec
}

func mkTestScreen(t *testing.T) tcell.SimulationScreen {
	t.Helper()
	s := tcell.NewSimulationScreen("")

	if s == nil {
		t.Fatalf("Failed to get simulation screen")
	}
	if e := s.Init(); e != nil {
		t.Fatalf("Failed to initialize screen: %v", e)
	}
	return s
}

func getPreviewContents(screen tcell.SimulationScreen) string {
	contents, w, h := screen.GetContents()

	startIndex := -1

	for i := range contents {
		runes := contents[i].Runes
		if len(runes) == 1 && runes[0] == '┌' {
			startIndex = i
		}
	}

	var indices []int

	prevLength := w - startIndex - 2

	for l := 1; l < h-1; l++ {
		for p := 0; p < prevLength; p++ {
			indices = append(indices, startIndex+1+l*w+p)
		}
	}

	result := ""
	for _, i := range indices {
		r := string(contents[i].Runes[0])
		result += r
	}

	return strings.TrimSpace(result)
}
