package uimsg

import "fmt"

func ConfigFileCreate(configPath string) string {
	return fmt.Sprintf(`Config file created at: %s

If you want to reset snipkit or delete the config, type in 'snipkit config clean'.`, configPath)
}

func ConfigFileDeleted(path string) string {
	return fmt.Sprintf(`Snipkit configuration file deleted: %s`, path)
}

func ThemesDeleted() string {
	return "Themes directory deleted"
}

func ThemesNotDeleted() string {
	return "Themes directory not deleted"
}

func ConfigNotDeleted() string {
	return "Config not deleted"
}

func ConfigNotFound(path string) string {
	return fmt.Sprintf(`No config found at: %s
You can set environment variable 'SNIPKIT_HOME' to change the directory path where snipkit will look for the config file.
Type in 'snipkit config init' to create a configuration file.
`, path)
}

func HomeDirectoryStillExists(path string) string {
	return fmt.Sprintf(`The snipkit home directory still exists exists since it holds non-deleted data (%s).
Please check for yourself if it can be deleted safely.`, path)
}

func ConfirmRecreateConfigFile(path string) string {
	return fmt.Sprintf("The configuration file already exists at %s.\nDo you want to recreate it?", path)
}

func ConfirmCreateConfigFile(path string) string {
	return fmt.Sprintf("Do you want to create a configuration file at %s?", path)
}

func ConfirmDeleteConfigFile(path string) string {
	return fmt.Sprintf("Do you really want to delete the snipkit configuration file (%s) ?", path)
}

func ConfirmDeleteThemesDir(path string) string {
	return fmt.Sprintf("The themes directory is not emtpty (%s). Should the custom themes be deleted as well?", path)
}
