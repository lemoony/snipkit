package uimsg

import (
	"bytes"
	"embed"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	configNotFound = "config_not_found.gotmpl"

	configFileCreateConfirm = "config_file_create_confirm.gotmpl"
	configFileCreateResult  = "config_file_create_result.gotmpl"

	configFileDeleteConfirm = "config_file_delete_confirm.gotmpl"
	configFileDeleteResult  = "config_file_delete_result.gotmpl"

	themesDeleteConfirm = "themes_delete_confirm.gotmpl"
	themesDeleteResult  = "themes_delete_result.gotmpl"

	homeDirStillExists = "home_dir_still_exists.gotmpl"

	managerAddConfigConfirm = "manager_add_config_confirm.gotmpl"
	managerAddConfigResult  = "manager_add_config_result.gotmpl"
)

var (
	highlightColor string
	colorProfile   = termenv.ColorProfile()

	snippetTextColor       = lipgloss.Color("#FFF7DB")
	snippetBackgroundColor = lipgloss.Color("#FB3082")

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

func ConfigFileCreateConfirm(path string, homeEnv string, recreate bool) Confirm {
	prompt := "Do you want to create the config file at this path?"
	if recreate {
		prompt = "Do you want to reset the config file?"
	}

	return Confirm{
		Prompt:   prompt,
		template: configFileCreateConfirm,
		data: map[string]interface{}{
			"homeEnvSet": homeEnv != "",
			"homeEnv":    homeEnv,
			"cfgPath":    path,
			"recreate":   recreate,
		},
	}
}

func ConfigFileCreateResult(created bool, configPath string, recreate bool) string {
	return render(
		configFileCreateResult,
		map[string]interface{}{
			"cfgPath":  configPath,
			"created":  created,
			"recreate": recreate,
		})
}

func ConfigFileDeleteConfirm(path string) Confirm {
	return Confirm{
		Prompt:   "Do you want to the config file?",
		template: configFileDeleteConfirm,
		data:     map[string]interface{}{"cfgPath": path},
	}
}

func ConfigFileDeleteResult(deleted bool, configPath string) string {
	return render(configFileDeleteResult, map[string]interface{}{"deleted": deleted, "cfgPath": configPath})
}

func ThemesDeleteConfirm(path string) Confirm {
	return Confirm{
		Prompt:   "Do you want to the delete the themes directory?",
		template: themesDeleteConfirm,
		data:     map[string]interface{}{"themesPath": path},
	}
}

func ThemesDeleteResult(deleted bool, themesPath string) string {
	return render(themesDeleteResult, map[string]interface{}{"deleted": deleted, "themesPath": themesPath})
}

func ManagerConfigAddConfirm(cfg string) Confirm {
	return Confirm{
		Prompt:   "Do you want to apply the change?",
		template: managerAddConfigConfirm,
		data:     map[string]interface{}{"configYaml": cfg},
	}
}

func ManagerAddConfigResult(confirmed bool, cfgPath string) string {
	return render(managerAddConfigResult, map[string]interface{}{"confirmed": confirmed, "cfgPath": cfgPath})
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
		"Snippet": func(values ...interface{}) string {
			blockStyle := lipgloss.NewStyle().
				Align(lipgloss.Left).
				Foreground(snippetTextColor).
				Background(snippetBackgroundColor).
				Padding(0).Margin(0)

			raw := strings.TrimSpace(values[0].(string))
			return blockStyle.Render(termenv.String(raw).String())
		},
		"Title": func(values ...interface{}) string {
			titleStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 1).
				SetString(values[0].(string))

			return titleStyle.String()
		},
	}
}
