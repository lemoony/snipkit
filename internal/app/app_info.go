package app

import (
	"fmt"

	"github.com/lemoony/snipkit/internal/model"
)

func (a *appImpl) Info() {
	a.printInfo(a.configService.Info())
	for _, manager := range a.managers {
		a.printInfo(manager.Info())
	}
}

func (a *appImpl) printInfo(info []model.InfoLine) {
	for _, line := range info {
		if line.IsError {
			a.tui.PrintError(fmt.Sprintf("%s: %s", line.Key, line.Value))
		} else {
			a.tui.PrintMessage(fmt.Sprintf("%s: %s", line.Key, line.Value))
		}
	}
}
