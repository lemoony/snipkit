package ui

import (
	"fmt"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lemoony/snippet-kit/internal/model"
)

var lexerMapping = map[model.Language]string{
	model.LanguageYAML:     "yaml",
	model.LanguageBash:     "bash",
	model.LanguageMarkdown: "markdown",
	model.LanguageTOML:     "toml",
}

func ShowLookup(snippets []model.Snippet) (int, error) {
	app := tview.NewApplication()

	selectedView := tview.NewTextView()
	selectedView.SetBorder(true)
	selectedView.SetTitle("Preview")
	selectedView.SetDynamicColors(true)

	selectedSnippet := -1

	previewWriter := tview.ANSIWriter(selectedView)
	f, s := getPreviewFormatterAndStyle()

	finder := tview.NewFinder().
		SetWrapAround(true).
		SetItems(len(snippets), func(index int) string {
			return snippets[index].Title
		}).
		SetDoneFunc(func(index int) {
			app.Stop()
			selectedSnippet = index
		}).
		SetChangedFunc(func(index int) {
			if index >= 0 {
				selectedView.SetText(fmt.Sprintf("Title: %s\n\n", snippets[index].Title))
				l := lexers.Get(lexerMapping[snippets[index].Language])
				if l == nil {
					l = lexers.Fallback
				}
				l = chroma.Coalesce(l)
				it, err := l.Tokenise(nil, snippets[index].Content)
				if err != nil {
					_, _ = selectedView.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				}
				err = f.Format(previewWriter, s, it)
				if err != nil {
					_, _ = selectedView.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				}
				selectedView.ScrollToBeginning()
			} else {
				selectedView.SetText("")
			}
		})

	finder.SetSelectedItemStyle(tcell.StyleDefault.Background(tview.Styles.ContrastBackgroundColor).Foreground(tview.Styles.PrimitiveBackgroundColor))
	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(finder, 0, 1, true).
		AddItem(selectedView, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		return -1, err
	}

	return selectedSnippet, nil
}

func getPreviewFormatterAndStyle() (chroma.Formatter, *chroma.Style) {
	f := formatters.Get("terminal")
	if f == nil {
		f = formatters.Fallback
	}

	s := styles.Get("friendly")
	if s == nil {
		s = styles.Fallback
	}

	return f, s
}
