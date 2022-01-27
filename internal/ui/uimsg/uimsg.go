package uimsg

import (
	"bytes"
	"embed"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
	"github.com/muesli/termenv"
	"golang.org/x/term"

	"github.com/lemoony/snipkit/internal/ui/style"
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

	snippetWidthMargin = 10
)

//go:embed templates/*.gotmpl
var templateFilesFS embed.FS

type Confirm struct {
	Prompt string

	header   string
	template string
	data     map[string]interface{}
}

type Printable struct {
	template string
	data     map[string]interface{}
}

func (p Printable) RenderWith(styler *style.Style) string {
	return renderWithStyle(p.template, styler, p.data)
}

func NewConfirm(prompt, header string) Confirm {
	return Confirm{
		Prompt: prompt,
		header: header,
	}
}

func (c *Confirm) Header(styler *style.Style, width int) string {
	if c.header != "" {
		return c.header
	}

	c.data["screenWidth"] = width

	return renderWithStyle(c.template, styler, c.data)
}

func (c *Confirm) HasTemplateHeader() bool {
	return c.template != ""
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

func ConfigFileCreateResult(created bool, configPath string, recreate bool) Printable {
	return Printable{
		template: configFileCreateResult,
		data: map[string]interface{}{
			"cfgPath":  configPath,
			"created":  created,
			"recreate": recreate,
		},
	}
}

func ConfigFileDeleteConfirm(path string) Confirm {
	return Confirm{
		Prompt:   "Do you want to delete the config file?",
		template: configFileDeleteConfirm,
		data:     map[string]interface{}{"cfgPath": path},
	}
}

func ConfigFileDeleteResult(deleted bool, configPath string) Printable {
	return Printable{
		template: configFileDeleteResult,
		data:     map[string]interface{}{"deleted": deleted, "cfgPath": configPath},
	}
}

func ThemesDeleteConfirm(path string) Confirm {
	return Confirm{
		Prompt:   "Do you want to the delete the themes directory?",
		template: themesDeleteConfirm,
		data:     map[string]interface{}{"themesPath": path},
	}
}

func ThemesDeleteResult(deleted bool, themesPath string) Printable {
	return Printable{
		template: themesDeleteResult,
		data:     map[string]interface{}{"deleted": deleted, "themesPath": themesPath},
	}
}

func ManagerConfigAddConfirm(cfg string) Confirm {
	return Confirm{
		Prompt:   "Do you want to apply the change?",
		template: managerAddConfigConfirm,
		data:     map[string]interface{}{"configYaml": cfg},
	}
}

func ManagerAddConfigResult(confirmed bool, cfgPath string) Printable {
	return Printable{
		template: managerAddConfigResult,
		data:     map[string]interface{}{"confirmed": confirmed, "cfgPath": cfgPath},
	}
}

func ConfigNotFound(configPath string) Printable {
	return Printable{
		template: configNotFound,
		data:     map[string]interface{}{"cfgPath": configPath},
	}
}

func HomeDirectoryStillExists(configPath string) Printable {
	return Printable{
		template: homeDirStillExists,
		data:     map[string]interface{}{"cfgPath": configPath},
	}
}

func renderWithStyle(templateFile string, styler *style.Style, data interface{}) string {
	t := newTemplate(templateFile, styler)
	writer := bytes.NewBufferString("")
	if err := t.Execute(writer, data); err != nil {
		panic(err)
	}
	return writer.String()
}

func newTemplate(fileName string, styler *style.Style) *template.Template {
	t, err := template.
		New(fileName).
		Funcs(termenv.TemplateFuncs(styler.ColorProfile())).
		Funcs(templateFuncs(styler)).
		ParseFS(templateFilesFS, filepath.Join("templates", fileName))
	if err != nil {
		panic(err)
	}
	return t
}

func templateFuncs(styler *style.Style) template.FuncMap {
	return template.FuncMap{
		"Highlighted": func(values ...interface{}) string {
			return lipgloss.
				NewStyle().
				Italic(true).
				Underline(true).
				Foreground(styler.HighlightColor().Value()).
				Render(values[0].(string))
		},
		"Snippet": func(values ...interface{}) string {
			width, _, _ := term.GetSize(0)
			width -= snippetWidthMargin

			blockStyle := lipgloss.NewStyle().
				Align(lipgloss.Left).
				Foreground(styler.SnippetContrastColor().Value()).
				Background(styler.SnippetColor().Value()).
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).BorderRight(true).BorderBottom(true).BorderLeft(true).
				BorderForeground(styler.SnippetColor().Value()).
				Padding(0).
				Margin(0).
				Width(width)

			raw := strings.TrimSpace(values[0].(string))
			raw = wrap.String(raw, width)

			return blockStyle.Render(raw)
		},
		"Title": func(values ...interface{}) string {
			return styler.Title(values[0].(string))
		},
	}
}
