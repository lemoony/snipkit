package ui

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/confirm"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/sync"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_Print(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("No config found at: /path/to/config")
	}, func(stdio termutil.Stdio) {
		term := NewTUI(WithStdio(stdio))
		term.Print(uimsg.ConfigNotFound("/path/to/config"))
		time.Sleep(time.Millisecond * 100)
	})
}

func Test_PrintMessage(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Hello world")
	}, func(stdio termutil.Stdio) {
		term := NewTUI(WithStdio(stdio))
		term.PrintMessage("Hello world")
		time.Sleep(time.Millisecond * 100)
	})
}

func Test_PrintError(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Some error message")
	}, func(stdio termutil.Stdio) {
		term := NewTUI(WithStdio(stdio))
		term.PrintError("Some error message")
		time.Sleep(time.Millisecond * 100)
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
		ttx := tt
		t.Run(tt.name, func(t *testing.T) {
			if ttx.envEditor == "" {
				_ = os.Unsetenv(envEditor)
			} else {
				_ = os.Setenv(envEditor, ttx.envEditor)
			}
			if ttx.envVisual == "" {
				_ = os.Unsetenv(envVisual)
			} else {
				_ = os.Setenv(envVisual, ttx.envVisual)
			}

			defer func() {
				_ = os.Unsetenv(envVisual)
				_ = os.Unsetenv(envEditor)
			}()

			editor := getEditor(ttx.preferred)
			assert.Equal(t, ttx.expected, editor)
		})
	}
}

func Test_ShowLookup(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{
			Title:    "Title 1",
			Content:  "Content: One",
			Language: model.LanguageYAML,
		},
		testutil.TestSnippet{
			Title:    "Title 2",
			Content:  "Content: Two",
			Language: model.LanguageYAML,
		},
	}

	runScreenTest(t, func(s tcell.Screen) {
		term := NewTUI(WithScreen(s))
		selected := term.ShowLookup(snippets, false)
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
	_ = os.Unsetenv("EDITOR")
	_ = os.Unsetenv("VISUAL")

	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.Send("iHello world\x1b")
		c.Send(":wq!\n")
	}, func(stdio termutil.Stdio) {
		testFile := path.Join(t.TempDir(), "testfile")
		_, err := os.Create(testFile)
		assert.NoError(t, err)

		NewTUI(WithStdio(stdio)).OpenEditor(testFile, "")
		bytes, err := os.ReadFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, "Hello world\n", string(bytes))
	})
}

func Test_OpenEditor_InvalidCommand(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		// nothing to expect since panic will be handled at application root level
	}, func(stdio termutil.Stdio) {
		testFile := path.Join(t.TempDir(), "testfile")
		_, err := os.Create(testFile)
		assert.NoError(t, err)

		assert.Panics(t, func() {
			NewTUI(WithStdio(stdio)).OpenEditor(testFile, "foo-editor")
		})
	})
}

func Test_Confirmation(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.Send("y")
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		term := NewTUI(WithStdio(stdio))
		confirmed := term.Confirmation(
			uimsg.ConfigFileDeleteConfirm("/some/path"),
			confirm.WithIn(stdio.In),
			confirm.WithOut(stdio.Out),
		)
		assert.True(t, confirmed)
	})
}

func Test_ShowSync(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Syncing all managers...")
		c.ExpectString("All done.")
	}, func(stdio termutil.Stdio) {
		term := NewTUI(WithStdio(stdio))
		screen := term.ShowSync()
		syncChannel := make(chan struct{})
		go func() {
			defer close(syncChannel)
			syncChannel <- struct{}{}
			screen.Start()
		}()
		<-syncChannel // wait for screen.Start()
		screen.Send(sync.UpdateStateMsg{Status: model.SyncStatusStarted})
		screen.Send(sync.UpdateStateMsg{Status: model.SyncStatusFinished})
		<-syncChannel // wait for screen.Start() to return
	})
}

func Test_ShowPicker(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Which snippet manager should be added to your configuration")
		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyUp)
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		term := NewTUI(WithStdio(stdio))
		index, ok := term.ShowPicker([]picker.Item{
			picker.NewItem("title1", "desc1"),
			picker.NewItem("title2", "desc2"),
			picker.NewItem("title3", "desc3"),
		})
		assert.Equal(t, 1, index)
		assert.True(t, ok)
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
