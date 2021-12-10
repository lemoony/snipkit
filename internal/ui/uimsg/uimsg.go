package uimsg

import "fmt"

func PrintNoConfig() {
	msg := `No snipkit configuration file found. Type in 'snipkit config init' to create one.`
	fmt.Println(msg)
}

func PrintConfigFileCreate(configPath string) {
	msg := `Config file created at: %s

If you want to reset snipkit or delete the config, type in 'snipkit config clean'.`

	fmt.Printf(msg+"\n", configPath)
}

func PrintConfigDeleted(path string) {
	msg := `Snipkit configuration file deleted: %s`
	fmt.Printf(msg+"\n", path)
}
