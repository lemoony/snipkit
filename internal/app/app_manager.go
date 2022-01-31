package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

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
	app := sync.Show()
	go a.startSyncManagers(app)
	_ = app.Start()
}

func (a *appImpl) startSyncManagers(p *tea.Program) {
	p.Send(sync.UpdateStateMsg{State: sync.State{Done: false}})

	for _, manager := range a.managers {
		sf := &model.SyncFeedback{
			Events: make(chan model.SyncEvent),
		}

		go func() {
			if !manager.Sync(sf) {
				close(sf.Events)
			}
		}()

		for v := range sf.Events {
			if v.State == model.SyncStateStarted {
				p.Send(sync.ManagerState{
					Key:        manager.Key(),
					InProgress: true,
					Error:      nil,
				})
			} else if v.State == model.SyncStateFinished {
				p.Send(sync.ManagerState{
					Key:        manager.Key(),
					InProgress: false,
					Error:      nil,
				})
			}
		}
	}

	p.Send(sync.UpdateStateMsg{State: sync.State{Done: true}})
}
