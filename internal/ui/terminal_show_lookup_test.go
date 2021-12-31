package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/model"
)

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

	screen := mkTestScreen(t)

	donec := make(chan struct{})
	go func() {
		defer close(donec)

		time.Sleep(time.Millisecond * 50)
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)))

		time.Sleep(time.Millisecond * 50)
		previewContent := getPreviewContents(screen)
		assert.Equal(t, snippets[1].Content, previewContent)

		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)))
	}()

	selected, err := showLookup(snippets, screen)

	assert.NoError(t, err)
	assert.Equal(t, 1, selected)
	<-donec
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
