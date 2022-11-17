package ui

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lemoony/snipkit/internal/model"
)

var lexerMapping = map[model.Language]string{
	model.LanguageYAML:     "yaml",
	model.LanguageBash:     "bash",
	model.LanguageMarkdown: "markdown",
	model.LanguageTOML:     "toml",
}

func (t *tuiImpl) ShowLookup(snippets []model.Snippet, fuzzySearch bool) int {
	app := tview.NewApplication()
	if t.screen != nil {
		app.SetScreen(t.screen)
	}

	preview := createPreview()
	previewWriter := tview.ANSIWriter(preview)

	selectedSnippet := -1
	f, s := t.getPreviewFormatterAndStyle()

	finder := tview.NewFinder().
		SetWrapAround(true).
		SetItems(len(snippets), func(index int) string {
			return snippets[index].GetTitle()
		}).
		SetDoneFunc(func(index int) {
			app.Stop()
			selectedSnippet = index
		}).
		SetChangedFunc(func(index int) {
			if index >= 0 {
				preview.SetText("")

				l := lexers.Get(lexerMapping[snippets[index].GetLanguage()])
				if l == nil {
					l = lexers.Fallback
				}
				l = chroma.Coalesce(l)
				it, err := l.Tokenise(nil, snippets[index].GetContent())
				if err != nil {
					_, _ = preview.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				}
				err = f.Format(previewWriter, s, it)
				if err != nil {
					_, _ = preview.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				}
				preview.ScrollToBeginning()
			} else {
				preview.SetText("")
			}
		})

	if fuzzySearch {
		finder.SetMatcherFunc(fuzzyMatcher)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(finder, 0, 1, true).
		AddItem(preview, 0, 1, false)

	t.applyStyle(finder, preview)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return selectedSnippet
}

func createPreview() *tview.TextView {
	result := tview.NewTextView()
	result.SetBorder(true)
	result.SetTitle("Preview")
	result.SetDynamicColors(true)
	result.SetBorderPadding(0, 0, 1, 0)
	return result
}

func (t *tuiImpl) getPreviewFormatterAndStyle() (chroma.Formatter, *chroma.Style) {
	f := formatters.Get("terminal")
	if f == nil {
		f = formatters.Fallback
	}

	s := styles.Get(t.styler.PreviewColorSchemeName())
	if s == nil {
		s = styles.Fallback
	}

	return f, s
}

func (t *tuiImpl) applyStyle(finder *tview.Finder, preview *tview.TextView) {
	finder.SetSelectedItemLabel(">")
	finder.SetInputLabel(">")
	finder.SetInputLabelStyle(tcell.StyleDefault.Foreground(t.styler.ActiveColor().CellValue()))

	finder.SetItemLabelPadding(1)

	finder.SetItemLabelStyle(tcell.StyleDefault.Background(t.styler.ActiveColor().CellValue()))
	finder.SetItemStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(t.styler.TextColor().CellValue()))

	finder.SetSelectedItemLabelStyle(tcell.StyleDefault.
		Background(t.styler.ActiveColor().CellValue()).
		Foreground(t.styler.ActiveContrastColor().CellValue()),
	)
	finder.SetSelectedItemStyle(tcell.StyleDefault.
		Background(t.styler.ActiveColor().CellValue()).
		Foreground(t.styler.ActiveContrastColor().CellValue()),
	)

	finder.SetCounterStyle(tcell.StyleDefault.Foreground(t.styler.InfoColor().CellValue()))
	finder.SetHighlightMatchStyle(tcell.StyleDefault.Foreground(t.styler.HighlightColor().CellValue()))
	finder.SetHighlightMatchMaintainBackgroundColor(true)

	finder.SetFieldStyle(tcell.StyleDefault.Foreground(t.styler.TextColor().CellValue()))
	finder.SetPlaceholderStyle(tcell.StyleDefault.Foreground(t.styler.PlaceholderColor().CellValue()))

	preview.SetTextColor(t.styler.TextColor().CellValue())
}

func fuzzyMatcher(slice string, input string) ([][2]int, int, bool) {
	slice = strings.TrimSpace(strings.ToLower(slice))
	input = strings.TrimSpace(strings.ToLower(input))

	fields := strings.Fields(input)
	score := 0
	var indices [][2]int

	for _, field := range fields {
		if found, subrange := partiallyContains(slice, field); found {
			s := subrange[1] - subrange[0]
			if s >= len(field) {
				score += s
				indices = append(indices, subrange)
			}
		}
	}

	return indices, score, score > 0
}

func partiallyContains(s string, substr string) (bool, [2]int) {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)

	for i := len(substr); i > 0; i-- {
		if index := strings.Index(s, substr[:i]); index >= 0 {
			return true, [2]int{index, index + i}
		}
	}

	return false, [2]int{}
}
