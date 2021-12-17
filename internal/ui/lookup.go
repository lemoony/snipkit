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

	preview := tview.NewTextView()
	preview.SetBorder(true)
	preview.SetTitle("Preview")
	preview.SetDynamicColors(true)
	preview.SetTextColor(currentTheme.previewSnippetNameColor())

	selectedSnippet := -1
	previewWriter := tview.ANSIWriter(preview)
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
				preview.SetText(fmt.Sprintf("Title: %s\n\n", snippets[index].Title))
				l := lexers.Get(lexerMapping[snippets[index].Language])
				if l == nil {
					l = lexers.Fallback
				}
				l = chroma.Coalesce(l)
				it, err := l.Tokenise(nil, snippets[index].Content)
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

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(finder, 0, 1, true).
		AddItem(preview, 0, 1, false)

	applyStyle(finder, preview, s)

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

	s := styles.Get(currentTheme.SyntaxHighlightingColorSchemeName)
	if s == nil {
		s = styles.Fallback
	}

	return f, s
}

func applyStyle(finder *tview.Finder, preview *tview.TextView, chromaStyle *chroma.Style) {
	finder.SetItemStyle(currentTheme.itemStyle())
	finder.SetSelectedItemStyle(currentTheme.selectedItemStyle())
	finder.SetCounterStyle(currentTheme.counterStyle())
	finder.SetHighlightMatchStyle(currentTheme.highlightItemMatchStyle())
	finder.SetFieldStyle(currentTheme.lookupInputStyle())
	finder.SetPlaceholderStyle(currentTheme.lookupInputPlaceholderStyle())

	if !currentTheme.SyntaxHighlightingApplyMainBackground {
		preview.SetBackgroundColor(tcell.GetColor(chromaStyle.Get(chroma.Background).Background.String()))
	}
}
