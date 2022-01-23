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

func (c *cliTerminal) ShowLookup(snippets []model.Snippet) int {
	app := tview.NewApplication()

	if c.screen != nil {
		app.SetScreen(c.screen)
	}

	preview := tview.NewTextView()
	preview.SetBorder(true)
	preview.SetTitle("Preview")
	preview.SetDynamicColors(true)
	preview.SetBorderPadding(0, 0, 1, 0)

	selectedSnippet := -1
	previewWriter := tview.ANSIWriter(preview)
	f, s := c.getPreviewFormatterAndStyle()

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

	c.applyStyle(finder, preview)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return selectedSnippet
}

func (c *cliTerminal) getPreviewFormatterAndStyle() (chroma.Formatter, *chroma.Style) {
	f := formatters.Get("terminal")
	if f == nil {
		f = formatters.Fallback
	}

	s := styles.Get(c.styler.PreviewColorSchemeName())
	if s == nil {
		s = styles.Fallback
	}

	return f, s
}

func (c *cliTerminal) applyStyle(finder *tview.Finder, preview *tview.TextView) {
	finder.SetSelectedItemLabel(">")
	finder.SetInputLabel(">")
	finder.SetInputLabelStyle(tcell.StyleDefault.Foreground(toColor(c.styler.ActiveColor())))

	finder.SetItemLabelPadding(1)

	finder.SetItemLabelStyle(tcell.StyleDefault.Background(toColor(c.styler.SelectionColor())))
	finder.SetItemStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(toColor(c.styler.TextColor())))

	finder.SetSelectedItemLabelStyle(tcell.StyleDefault.
		Background(toColor(c.styler.SelectionColor())).
		Foreground(toColor(c.styler.SelectionColorReverse())),
	)
	finder.SetSelectedItemStyle(tcell.StyleDefault.
		Background(toColor(c.styler.SelectionColor())).
		Foreground(toColor(c.styler.SelectionColorReverse())),
	)

	finder.SetCounterStyle(tcell.StyleDefault.Foreground(toColor(c.styler.InfoColor())))
	finder.SetHighlightMatchStyle(tcell.StyleDefault.Foreground(toColor(c.styler.HighlightColor())))
	finder.SetHighlightMatchMaintainBackgroundColor(true)

	finder.SetFieldStyle(tcell.StyleDefault.Foreground(toColor(c.styler.TextColor())))
	finder.SetPlaceholderStyle(tcell.StyleDefault.Foreground(toColor(c.styler.PlaceholderColor())))

	preview.SetTextColor(toColor(c.styler.TextColor()))
}
