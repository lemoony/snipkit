package app

import (
	"strings"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/sync"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
)

func (a *appImpl) AddManager() {
	managerDescriptions := a.provider.ManagerDescriptions(a.config.Manager)
	listItems := make([]picker.Item, len(managerDescriptions))
	for i := range managerDescriptions {
		listItems[i] = picker.NewItem(managerDescriptions[i].Name, managerDescriptions[i].Description)
	}

	if index, ok := a.tui.ShowPicker(listItems); ok {
		managerDescription := managerDescriptions[index]
		cfg := a.provider.AutoConfig(managerDescription.Key, a.system)
		configBytes := config.SerializeToYamlWithComment(cfg)
		configStr := strings.TrimSpace(string(configBytes))
		confirmed := a.tui.Confirmation(uimsg.ManagerConfigAddConfirm(configStr))
		if confirmed {
			a.configService.UpdateManagerConfig(cfg)
		}
		a.tui.Print(uimsg.ManagerAddConfigResult(confirmed, a.configService.ConfigFilePath()))
	}
}

func (a *appImpl) SyncManager() {
	screen := sync.New()
	go a.startSyncManagers(screen)
	screen.Start()
}

func (a *appImpl) startSyncManagers(s *sync.Screen) {
	s.Send(sync.UpdateStateMsg{Status: model.SyncStatusStarted})

	for _, manager := range a.managers {
		events := make(chan model.SyncEvent)

		go func() {
			if !manager.Sync(events) {
				close(events)
			}
		}()

		for v := range events {
			s.Send(sync.UpdateStateMsg{
				ManagerState: &sync.ManagerState{
					Key:    manager.Key(),
					Status: v.Status,
					Lines:  v.Lines,
					Input:  v.Login,
				},
			})
		}
	}

	s.Send(sync.UpdateStateMsg{Status: model.SyncStatusStarted})
}
