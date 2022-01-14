package uimsg

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"text/template"
)

const (
	configFileCreated             = "config_file_created.gotmpl"
	configFileDelete              = "config_file_deleted.gotmpl"
	configNotFound                = "config_not_found.gotmpl"
	configFileCreateDescription   = "config_file_create_description.gotmpl"
	configFileRecreateDescription = "config_file_recreate_description.gotmpl"
	homeDirStillExists            = "home_dir_still_exists.gotmpl"
	themesDirDeleteConfirm        = "themes_dir_delete_confirm.gotmpl"
)

//go:embed templates/*.gotmpl
var templateFilesFS embed.FS

func ConfigFileCreated(configPath string) string {
	return render(configFileCreated, map[string]interface{}{"cfgPath": configPath})
}

func ConfigFileDeleted(configPath string) string {
	return render(configFileDelete, map[string]interface{}{"cfgPath": configPath})
}

func ConfigNotFound(configPath string) string {
	return render(configNotFound, map[string]interface{}{"cfgPath": configPath})
}

func ConfigNotDeleted() string {
	return "Config not deleted"
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

func ConfigFileCreateConfirm(path string) string {
	return fmt.Sprintf(`Do you want to create a configuration file at %s?`, path)
}

func ConfigFileDeleteConfirm(path string) string {
	return fmt.Sprintf("Do you really want to delete the snipkit configuration file at %s?", path)
}

func ThemesDirDeleteConfirm(path string) string {
	return render(themesDirDeleteConfirm, map[string]interface{}{"themesPath": path})
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
	t, err := template.ParseFS(templateFilesFS, filepath.Join("templates", fileName))
	if err != nil {
		panic(err)
	}
	return t
}
