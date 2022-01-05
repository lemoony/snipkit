package app

import (
	"fmt"
)

func (a *appImpl) Info() {
	a.ui.PrintMessage(fmt.Sprintf("%s: %s", "Config file", a.configService.ConfigFilePath()))

	for _, provider := range a.Providers {
		for _, line := range provider.Info().Lines {
			if line.IsError {
				a.ui.PrintError(fmt.Sprintf("%s: %s", line.Key, line.Value))
			} else {
				a.ui.PrintMessage(fmt.Sprintf("%s: %s", line.Key, line.Value))
			}
		}
	}
}
