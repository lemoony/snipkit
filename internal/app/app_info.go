package app

import (
	"fmt"

	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

func (a *appImpl) Info() {
	a.tui.PrintMessage(fmt.Sprintf("%s: %s", "Config file", a.configService.ConfigFilePath()))
	a.tui.PrintMessage(fmt.Sprintf("%s: %s", "SNIPKIT_HOME",
		stringutil.StringOrDefault(a.system.HomeEnvValue(), "Not set")),
	)

	for _, manager := range a.managers {
		for _, line := range manager.Info().Lines {
			if line.IsError {
				a.tui.PrintError(fmt.Sprintf("%s: %s", line.Key, line.Value))
			} else {
				a.tui.PrintMessage(fmt.Sprintf("%s: %s", line.Key, line.Value))
			}
		}
	}
}
