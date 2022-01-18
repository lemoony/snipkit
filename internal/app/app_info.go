package app

import (
	"fmt"

	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

func (a *appImpl) Info() {
	a.ui.PrintMessage(fmt.Sprintf("%s: %s", "Config file", a.configService.ConfigFilePath()))
	a.ui.PrintMessage(fmt.Sprintf("%s: %s", "SNIPKIT_HOME",
		stringutil.StringOrDefault(a.system.HomeEnvValue(), "Not set")),
	)

	for _, manager := range a.managers {
		for _, line := range manager.Info().Lines {
			if line.IsError {
				a.ui.PrintError(fmt.Sprintf("%s: %s", line.Key, line.Value))
			} else {
				a.ui.PrintMessage(fmt.Sprintf("%s: %s", line.Key, line.Value))
			}
		}
	}
}
