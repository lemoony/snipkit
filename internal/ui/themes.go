package ui

import (
	"github.com/gdamore/tcell/v2"
)

type simpleThemeValues struct {
	syntaxHighlightingColorSchemeName     string
	syntaxHighlightingApplyMainBackground bool
	backgroundColor                       string
	itemTextColor                         string
	selectedItemTextColor                 string
	selectionColor                        string
	highlightColor                        string
	placeholderColor                      string
	secondaryLabelColor                   string
}

func (s *simpleThemeValues) toThemeValues() ThemeValues {
	return ThemeValues{
		SyntaxHighlightingColorSchemeName:     s.syntaxHighlightingColorSchemeName,
		SyntaxHighlightingApplyMainBackground: s.syntaxHighlightingApplyMainBackground,
		PreviewSnippetDefaultTextColor:        "",
		BackgroundColor:                       s.backgroundColor,
		BorderColor:                           s.itemTextColor,
		BorderTitleColor:                      s.itemTextColor,
		ItemTextColor:                         s.itemTextColor,
		SelectedItemTextColor:                 s.selectedItemTextColor,
		SelectedItemBackgroundColor:           s.selectionColor,
		ItemHighlightMatchColor:               s.highlightColor,
		CounterTextColor:                      s.secondaryLabelColor,
		CounterBackgroundColor:                "",
		LookupInputTextColor:                  s.itemTextColor,
		LookupInputPlaceholderColor:           s.placeholderColor,
		LookupInputBackgroundColor:            "",
		LookupLabelTextColor:                  s.secondaryLabelColor,
	}
}

var themeNames = map[string]ThemeValues{
	"default": themeDefault,
	"dracula": themeDracula.toThemeValues(),
}

var themeDefault = themeSimple.toThemeValues()

var themeSimple = simpleThemeValues{
	syntaxHighlightingColorSchemeName:     "friendly",
	syntaxHighlightingApplyMainBackground: true,
	backgroundColor:                       "",
	itemTextColor:                         "",
	selectedItemTextColor:                 "#f8f8f2",
	selectionColor:                        "#ff5555",
	highlightColor:                        "#50fa7b",
	secondaryLabelColor:                   "#B79103",
	placeholderColor:                      "#6272a4",
}

var themeDracula = simpleThemeValues{
	syntaxHighlightingColorSchemeName:     "dracula",
	syntaxHighlightingApplyMainBackground: true,
	backgroundColor:                       "#282a36",
	itemTextColor:                         "#f8f8f2",
	selectedItemTextColor:                 "#f8f8f2",
	selectionColor:                        "#ff5555",
	highlightColor:                        "#50fa7b",
	secondaryLabelColor:                   "#f1fa8c",
	placeholderColor:                      "#6272a4",
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
	return tcell.StyleDefault.Background(r.backgroundColor()).Foreground(tcell.GetColor(r.ItemTextColor))
}

func (r *ThemeValues) itemLabelStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.SelectedItemBackgroundColor)).Foreground(tcell.GetColor(r.ItemTextColor))
}

func (r *ThemeValues) selectedItemLabelStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.SelectedItemBackgroundColor)).Foreground(tcell.GetColor(r.SelectedItemTextColor))
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

func (r *ThemeValues) lookupLabelStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.BackgroundColor)).Foreground(tcell.GetColor(r.LookupLabelTextColor))
}

func (r *ThemeValues) lookupInputPlaceholderStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.GetColor(r.LookupInputPlaceholderColor))
}

func (r *ThemeValues) previewSnippetDefaultTextColor() tcell.Color {
	return tcell.GetColor(r.PreviewSnippetDefaultTextColor)
}
