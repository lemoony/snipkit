package app

import (
	"strings"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/ui/confirm"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
)

func (a *appImpl) AddManager() {
	managerDescriptions := a.provider.ManagerDescriptions(a.config.Manager)
	listItems := make([]picker.Item, len(managerDescriptions))
	for i := range managerDescriptions {
		listItems[i] = picker.NewItem(managerDescriptions[i].Name, managerDescriptions[i].Description)
	}

	if index, ok := a.ui.ShowPicker(listItems); ok {
		managerDescription := managerDescriptions[index]
		cfg := a.provider.AutoConfig(managerDescription.Key, a.system)
		configBytes := config.SerializeToYamlWithComment(cfg)
		configStr := strings.TrimSpace(string(configBytes))
		confirmed := a.ui.Confirmation(uimsg.ManagerConfigAddConfirm(configStr), confirm.WithFullscreen())
		if confirmed {
			a.configService.UpdateManagerConfig(cfg)
		}
		a.ui.PrintMessage(uimsg.ManagerAddConfigResult(confirmed, a.configService.ConfigFilePath()))
	}
}
