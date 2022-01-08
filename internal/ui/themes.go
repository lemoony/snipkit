package ui

import (
	_ "embed"
	"path/filepath"
	"reflect"
	"regexp"

	"emperror.dev/errors"
	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"

	themedata "github.com/lemoony/snippet-kit/themes"
)

const (
	defaultThemeName       = "default"
	variablePatternMatches = 2
	filenamePatternMatches = 2
)

var (
	variablePattern = regexp.MustCompile(`^\${(?P<varName>.*)}$`)
	filenamePattern = regexp.MustCompile(`^(?P<filename>.*)\.ya?ml$`)
	ErrInvalidTheme = errors.New("invalid theme")
)

type themeWrapper struct {
	Version   string `yaml:"version"`
	Variables map[string]string
	Theme     ThemeValues `yaml:"theme"`
}

func embeddedTheme(name string) ThemeValues {
	entries, err := themedata.Files.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		m := filenamePattern.FindStringSubmatch(filepath.Base(entry.Name()))
		if len(m) == filenamePatternMatches {
			themeName := m[1]
			if name == themeName {
				return readTheme(entry.Name())
			}
		}
	}

	panic(errors.Wrapf(ErrInvalidTheme, "theme not found: "+name))
}

func readTheme(path string) ThemeValues {
	bytes, err := themedata.Files.ReadFile(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read theme %s", path))
	}

	var wrapper themeWrapper
	err = yaml.Unmarshal(bytes, &wrapper)
	if err != nil {
		panic(err)
	}

	return wrapper.theme()
}

func (t *themeWrapper) theme() ThemeValues {
	result := t.Theme
	v := reflect.Indirect(reflect.ValueOf(&result))

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.String {
			matches := variablePattern.FindStringSubmatch(v.Field(i).String())
			if len(matches) != variablePatternMatches {
				continue
			}
			varName := matches[1]
			if val, ok := t.Variables[varName]; !ok {
				panic(errors.Wrapf(ErrInvalidTheme, "variable %s not found", varName))
			} else {
				v.Field(i).SetString(val)
			}
		}
	}
	return result
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
	return tcell.StyleDefault.Background(tcell.GetColor(r.ItemHighlightMatchBackgroundColor)).Foreground(tcell.GetColor(r.ItemHighlightMatchTextColor))
}

func (r *ThemeValues) counterStyle() tcell.Style {
	return tcell.StyleDefault.Foreground(tcell.GetColor(r.CounterTextColor))
}

func (r *ThemeValues) lookupInputStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.LookupInputBackgroundColor)).Foreground(tcell.GetColor(r.LookupInputTextColor))
}

func (r *ThemeValues) lookupLabelStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.GetColor(r.BackgroundColor))
}

func (r *ThemeValues) lookupInputPlaceholderStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.GetColor(r.LookupInputPlaceholderColor))
}

func (r *ThemeValues) parametersLabelColor() tcell.Color {
	return tcell.GetColor(r.ParametersLabelTextColor)
}

func (r *ThemeValues) parametersFieldBackgroundColor() tcell.Color {
	return tcell.GetColor(r.ParametersFieldBackgroundColor)
}

func (r *ThemeValues) parametersFieldTextColor() tcell.Color {
	return tcell.GetColor(r.ParametersFieldTextColor)
}

func (r *ThemeValues) parametersAutocompleteTextColor() tcell.Color {
	return tcell.GetColor(r.ParameterAutocompleteTextColor)
}

func (r *ThemeValues) parametersAutocompleteSelectedTextColor() tcell.Color {
	return tcell.GetColor(r.ParameterAutocompleteSelectedTextColor)
}

func (r *ThemeValues) parametersAutocompleteBackgroundColor() tcell.Color {
	return tcell.GetColor(r.ParameterAutocompleteBackgroundColor)
}

func (r *ThemeValues) parametersAutocompleteSelectedBackgroundColor() tcell.Color {
	return tcell.GetColor(r.ParameterAutocompleteSelectedBackgroundColor)
}

func (r *ThemeValues) selectedButtonBackgroundColor() tcell.Color {
	return tcell.GetColor(r.SelectedItemBackgroundColor)
}

func (r *ThemeValues) previewDefaultTextColor() tcell.Color {
	return tcell.GetColor(r.PreviewDefaultTextColor)
}

func (r *ThemeValues) previewOverwriteBackgroundColor() (tcell.Color, bool) {
	if r.PreviewOverwriteBackgroundColor != "" {
		return tcell.GetColor(r.PreviewOverwriteBackgroundColor), true
	} else {
		return tcell.ColorDefault, false
	}
}
