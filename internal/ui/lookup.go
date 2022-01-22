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

func (c cliTerminal) ShowLookup(snippets []model.Snippet) int {
	app := tview.NewApplication()

	if c.screen != nil {
		app.SetScreen(c.screen)
	}

	preview := tview.NewTextView()
	preview.SetBorder(true)
	preview.SetTitle("Preview")
	preview.SetDynamicColors(true)

	selectedSnippet := -1
	previewWriter := tview.ANSIWriter(preview)
	f, s := getPreviewFormatterAndStyle()

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

	applyStyle(finder, preview, s)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return selectedSnippet
}

func getPreviewFormatterAndStyle() (chroma.Formatter, *chroma.Style) {
	f := formatters.Get("terminal")
	if f == nil {
		f = formatters.Fallback
	}

	s := styles.Get(currentTheme.PreviewColorSchemeName)
	if s == nil {
		s = styles.Fallback
	}

	return f, s
}

func applyStyle(finder *tview.Finder, preview *tview.TextView, chromaStyle *chroma.Style) {
	finder.SetSelectedItemLabel(">")
	finder.SetInputLabel(">")
	finder.SetInputLabelStyle(currentTheme.lookupLabelStyle())

	finder.SetItemLabelPadding(1)

	finder.SetItemLabelStyle(tcell.StyleDefault.Background(tcellColor(styler.SelectionColor())))
	finder.SetItemStyle(currentTheme.itemStyle())

	finder.SetSelectedItemLabelStyle(tcell.StyleDefault.Background(tcellColor(styler.SelectionColor())).Foreground(tcellColor(styler.SelectionColorReverse())))
	finder.SetSelectedItemStyle(tcell.StyleDefault.Background(tcellColor(styler.SelectionColor())).Foreground(tcellColor(styler.SelectionColorReverse())))

	finder.SetCounterStyle(tcell.StyleDefault.Foreground(tcellColor(styler.InfoColor())))
	finder.SetHighlightMatchStyle(tcell.StyleDefault.Foreground(tcellColor(styler.HighlightColor())))
	finder.SetHighlightMatchMaintainBackgroundColor(true)

	finder.SetFieldStyle(currentTheme.lookupInputStyle())
	finder.SetPlaceholderStyle(currentTheme.lookupInputPlaceholderStyle())

	preview.SetTextColor(currentTheme.previewDefaultTextColor())
	if !currentTheme.PreviewApplyMainBackground {
		if bgColor, ok := currentTheme.previewOverwriteBackgroundColor(); ok {
			preview.SetBackgroundColor(bgColor)
		} else {
			preview.SetBackgroundColor(tcell.GetColor(chromaStyle.Get(chroma.Background).Background.String()))
		}
	}
}
