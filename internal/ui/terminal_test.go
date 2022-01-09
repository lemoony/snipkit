package ui

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/utils/testutil"
)

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
		confirmed := term.Confirm("Are you sure?")
		assert.True(t, confirmed)
	})
}

func Test_getEditor(t *testing.T) {
	tests := []struct {
		name      string
		envVisual string
		envEditor string
		preferred string
		expected  string
	}{
		{name: "default editor unix", expected: defaultEditor},
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
		})
	}
}

func Test_ShowLookup(t *testing.T) {
	snippets := []model.Snippet{
		{
			TitleFunc:    testutil.FixedString("Title 1"),
			ContentFunc:  testutil.FixedString("Content: One"),
			LanguageFunc: testutil.FixedLanguage(model.LanguageYAML),
		},
		{
			TitleFunc:    testutil.FixedString("Title 2"),
			ContentFunc:  testutil.FixedString("Content: Two"),
			LanguageFunc: testutil.FixedLanguage(model.LanguageYAML),
		},
	}

	runScreenTest(t, func(s tcell.Screen) {
		term := NewTerminal(WithScreen(s))
		selected := term.ShowLookup(snippets)
		assert.Equal(t, 1, selected)
	}, func(screen tcell.SimulationScreen) {
		time.Sleep(time.Millisecond * 50)
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)))

		time.Sleep(time.Millisecond * 50)
		previewContent := getPreviewContents(screen)
		assert.Equal(t, snippets[1].GetContent(), previewContent)

		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))
	})
}

func Test_OpenEditor(t *testing.T) {
	runExpectTest(t, func(c *expect.Console) {
		_, _ = c.Send("iHello world\x1b")
		time.Sleep(time.Second)
		_, _ = c.SendLine(":wq!")
	}, func(stdio terminal.Stdio) {
		term := NewTerminal(WithStdio(stdio))

		testFile := path.Join(t.TempDir(), "testfile")
		_, err := os.Create(testFile)
		assert.NoError(t, err)

		term.OpenEditor(testFile, "")
		bytes, err := ioutil.ReadFile(testFile) //nolint:gosec // potential file inclusion via variable
		assert.NoError(t, err)
		assert.Equal(t, "Hello world\n", string(bytes))
	})
}

func Test_OpenEditor_InvalidCommand(t *testing.T) {
	runExpectTest(t, func(c *expect.Console) {
		// nothing to expect since panic will be handled at application root level
	}, func(stdio terminal.Stdio) {
		term := NewTerminal(WithStdio(stdio))

		testFile := path.Join(t.TempDir(), "testfile")
		_, err := os.Create(testFile)
		assert.NoError(t, err)

		assert.Panics(t, func() {
			term.OpenEditor(testFile, "foo-editor")
		})
	})
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
