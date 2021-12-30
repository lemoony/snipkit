package uimsg

import "fmt"

func NoConfig() string {
	return `No snipkit configuration file found. Type in 'snipkit config init' to create one.`
}

func ConfigFileCreate(configPath string) string {
	return fmt.Sprintf(`Config file created at: %s

If you want to reset snipkit or delete the config, type in 'snipkit config clean'.`, configPath)
}

func ConfigFileDeleted(path string) string {
	return fmt.Sprintf(`Snipkit configuration file deleted: %s`, path)
}

func ConfirmRecreateConfigFile(path string) string {
	return fmt.Sprintf("The configuration file already exists at %s.\nDo you want to recreate it?", path)
}

func ConfirmCreateConfigFile() string {
	return "There is no snipkit config file currently. Do you want to create one?"
}

func ConfirmDeleteConfigFile() string {
	return "Do you really want to delete the snipkit configuration file?"
}
