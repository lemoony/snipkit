package uimsg

import (
	"bytes"
	"embed"
	"path/filepath"
	"text/template"

	"github.com/muesli/termenv"
)

const (
	configFileCreated             = "config_file_created.gotmpl"
	configFileDelete              = "config_file_deleted.gotmpl"
	configNotFound                = "config_not_found.gotmpl"
	configFileCreateDescription   = "config_file_create_description.gotmpl"
	configFileRecreateDescription = "config_file_recreate_description.gotmpl"
	configFileDeleteDescription   = "config_file_delete_description.gotmpl"
	themesDirDeleteDescription    = "themes_dir_delete_description.gotmpl"
	homeDirStillExists            = "home_dir_still_exists.gotmpl"
)

var (
	highlightColor string
	colorProfile   = termenv.ColorProfile()

	//go:embed templates/*.gotmpl
	templateFilesFS embed.FS
)

func ConfigFileCreated(configPath string) string {
	return render(configFileCreated, map[string]interface{}{"cfgPath": configPath})
}

func ConfigFileDeleted(configPath string) string {
	return render(configFileDelete, map[string]interface{}{"cfgPath": configPath})
}

func ConfigNotFound(configPath string) string {
	return render(configNotFound, map[string]interface{}{"cfgPath": configPath})
}

func ThemesDeleted() string {
	return "Themes directory deleted"
}

func ThemesNotDeleted() string {
	return "Themes directory not deleted"
}

func HomeDirectoryStillExists(configPath string) string {
	return render(homeDirStillExists, map[string]interface{}{"cfgPath": configPath})
}

func ConfigFileRecreateDescription(configPath string) string {
	return render(configFileRecreateDescription, map[string]interface{}{"cfgPath": configPath})
}

func ConfigFileDeleteDescription(configPath string) string {
	return render(configFileDeleteDescription, map[string]interface{}{"cfgPath": configPath})
}

func ThemesDirDeleteDescription(configPath string) string {
	return render(themesDirDeleteDescription, map[string]interface{}{"themesPath": configPath})
}

func ConfigFileRecreateConfirm() string {
	return "Do you want to recreate the config file?"
}

func ConfigFileCreateDescription(path string, homeEnv string) string {
	return render(configFileCreateDescription, map[string]interface{}{
		"homeEnvSet": homeEnv != "",
		"homeEnv":    homeEnv,
		"cfgPath":    path,
	})
}

func ConfigFileCreateConfirm() string {
	return "Do you want to create the config file at this path?"
}

func ConfigFileDeleteConfirm() string {
	return "Do you want to delete config file?"
}

func ThemesDirDeleteConfirm() string {
	return "Do you want to the custom themes as well?"
}

func render(templateFile string, data interface{}) string {
	t := newTemplate(templateFile)
	writer := bytes.NewBufferString("")
	if err := t.Execute(writer, data); err != nil {
		panic(err)
	}
	return writer.String()
}

func newTemplate(fileName string) *template.Template {
	t, err := template.
		New(fileName).
		Funcs(termenv.TemplateFuncs(colorProfile)).
		Funcs(templateFuncs()).
		ParseFS(templateFilesFS, filepath.Join("templates", fileName))
	if err != nil {
		panic(err)
	}
	return t
}

func SetHighlightColor(color string) {
	highlightColor = color
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"Highlighted": func(values ...interface{}) string {
			s := termenv.String(values[0].(string))
			s = s.Foreground(colorProfile.Color(highlightColor))
			s = s.Italic()
			return s.String()
		},
	}
}
