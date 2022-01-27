package ui

import (
	"fmt"

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

func (t *tuiImpl) ShowLookup(snippets []model.Snippet) int {
	app := tview.NewApplication()

	if t.screen != nil {
		app.SetScreen(t.screen)
	}

	preview := tview.NewTextView()
	preview.SetBorder(true)
	preview.SetTitle("Preview")
	preview.SetDynamicColors(true)
	preview.SetBorderPadding(0, 0, 1, 0)

	selectedSnippet := -1
	previewWriter := tview.ANSIWriter(preview)
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

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(finder, 0, 1, true).
		AddItem(preview, 0, 1, false)

	t.applyStyle(finder, preview)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return selectedSnippet
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
