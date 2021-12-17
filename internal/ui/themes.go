package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/phuslu/log"
	"github.com/rivo/tview"
)

var themeNames = map[string]ThemeValues{
	"default": themeDefault,
	"dracula": themeDracula,
}

var themeDefault = ThemeValues{
	SyntaxHighlightingColorSchemeName:     "friendly",
	SyntaxHighlightingApplyMainBackground: false,
	BackgroundColor:                       "0",
	BorderColor:                           "lightgray",
	BorderTitleColor:                      "",
	ItemTextColor:                         "",
	SelectedItemTextColor:                 "",
	SelectedItemBackgroundColor:           "red",
	ItemHighlightMatchColor:               "green",
	CounterTextColor:                      "yellow",
	CounterBackgroundColor:                "",
	LookupInputTextColor:                  "",
	LookupInputPlaceholderColor:           "lightgray",
	LookupInputBackgroundColor:            "",
	PreviewSnippetNameColor:               "",
}

var themeDracula = ThemeValues{
	SyntaxHighlightingColorSchemeName:     "dracula",
	SyntaxHighlightingApplyMainBackground: false,
	PreviewSnippetNameColor:               "white",
	BackgroundColor:                       "#00005F",
	BorderColor:                           "lightgray",
	BorderTitleColor:                      "white",
	ItemTextColor:                         "white",
	SelectedItemTextColor:                 "white",
	SelectedItemBackgroundColor:           "red",
	ItemHighlightMatchColor:               "green",
	CounterTextColor:                      "yellow",
	CounterBackgroundColor:                "",
	LookupInputTextColor:                  "white",
	LookupInputPlaceholderColor:           "lightgray",
	LookupInputBackgroundColor:            "",
}

var currentTheme ThemeValues

func ApplyConfig(cfg Config) {
	theme := themeDefault
	for _, t := range cfg.CustomThemes {
		if t.Name == cfg.Theme {
			log.Debug().Msgf("Applied custom theme: %s", t.Name)
			theme = t.Values
		}
	}

	if t, ok := themeNames[cfg.Theme]; ok {
		log.Debug().Msgf("Applied provided theme: %s", cfg.Theme)
		theme = t
	}

	setTheme(theme)
}

func setTheme(theme ThemeValues) {
	currentTheme = theme
	tview.Styles.PrimitiveBackgroundColor = theme.backgroundColor()
	tview.Styles.BorderColor = theme.borderColor()
	tview.Styles.TitleColor = theme.borderTitleColor()
}

func (r *ThemeValues) backgroundColor() tcell.Color {
	return tcell.GetColor(r.BackgroundColor)
}

func (r *ThemeValues) borderColor() tcell.Color {
	return tcell.GetColor(r.BorderColor)
}

func (r *ThemeValues) borderTitleColor() tcell.Color {
	return tcell.GetColor(r.BorderTitleColor)
}

func (r *ThemeValues) itemStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.GetColor(r.ItemTextColor))
}

func (r *ThemeValues) selectedItemStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.SelectedItemBackgroundColor)).Foreground(tcell.GetColor(r.SelectedItemTextColor))
}

func (r *ThemeValues) highlightItemMatchStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.GetColor(r.ItemHighlightMatchColor))
}

func (r *ThemeValues) counterStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.CounterBackgroundColor)).Foreground(tcell.GetColor(r.CounterTextColor))
}

func (r *ThemeValues) lookupInputStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.LookupInputBackgroundColor)).Foreground(tcell.GetColor(r.LookupInputTextColor))
}

func (r *ThemeValues) lookupInputPlaceholderStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.GetColor(r.LookupInputPlaceholderColor))
}

func (r *ThemeValues) previewSnippetNameColor() tcell.Color {
	return tcell.GetColor(r.PreviewSnippetNameColor)
}
