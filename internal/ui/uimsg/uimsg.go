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

type Confirm struct {
	Prompt string

	template string
	data     map[string]interface{}
}

func (c *Confirm) Header() string {
	return render(c.template, c.data)
}

func SetHighlightColor(color string) {
	highlightColor = color
}

func ConfirmConfigCreation(path string, homeEnv string) Confirm {
	return Confirm{
		Prompt:   "Do you want to create the config file at this path?",
		template: configFileCreateDescription,
		data: map[string]interface{}{
			"homeEnvSet": homeEnv != "",
			"homeEnv":    homeEnv,
			"cfgPath":    path,
		},
	}
}

func ConfirmConfigRecreate(path string) Confirm {
	return Confirm{
		Prompt:   "Do you want to reset the config file?",
		template: configFileRecreateDescription,
		data:     map[string]interface{}{"cfgPath": path},
	}
}

func ConfirmConfigDelete(path string) Confirm {
	return Confirm{
		Prompt:   "Do you want to delete the config file?",
		template: configFileDeleteDescription,
		data:     map[string]interface{}{"cfgPath": path},
	}
}

func ConfirmThemesDelete(path string) Confirm {
	return Confirm{
		Prompt:   "Do you want to the custom themes?",
		template: themesDirDeleteDescription,
		data:     map[string]interface{}{"themesPath": path},
	}
}

func ConfigFileCreated(configPath string) string {
	return render(configFileCreated, map[string]interface{}{"cfgPath": configPath})
}

func ConfigFileDeleted(configPath string) string {
	return render(configFileDelete, map[string]interface{}{"cfgPath": configPath})
}

func ConfigNotFound(configPath string) string {
	return render(configNotFound, map[string]interface{}{"cfgPath": configPath})
}

func HomeDirectoryStillExists(configPath string) string {
	return render(homeDirStillExists, map[string]interface{}{"cfgPath": configPath})
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
