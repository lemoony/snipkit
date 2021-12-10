package uisurvey

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func ConfirmRecreateConfigFile(path string) (bool, error) {
	create := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("The configuration file already exists at %s.\nDo you want to recreate it?", path),
	}

	if err := survey.AskOne(prompt, &create); err != nil {
		return false, err
	}

	return create, nil
}

func ConfirmCreateConfigFile() (bool, error) {
	create := false
	prompt := &survey.Confirm{
		Message: "There is no snipkit config file currently. Do you want to create one?",
	}

	if err := survey.AskOne(prompt, &create); err != nil {
		return false, err
	}

	return create, nil
}

func ConfirmDeleteConfigFile() (bool, error) {
	create := false
	prompt := &survey.Confirm{
		Message: "Do you really want to delete the snipkit configuration file?",
	}

	if err := survey.AskOne(prompt, &create); err != nil {
		return false, err
	}

	return create, nil
}
