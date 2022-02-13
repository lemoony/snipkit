package app

import (
	"strings"

	"emperror.dev/errors"
	"github.com/phuslu/log"

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
	syncScreen := a.tui.ShowSync()

	doneChannel := make(chan struct{})

	go func() {
		defer func() {
			close(doneChannel)
		}()
		a.startSyncManagers(syncScreen)
	}()

	syncScreen.Start()

	<-doneChannel
}

func (a *appImpl) startSyncManagers(syncScreen sync.Screen) {
	syncScreen.Send(sync.UpdateStateMsg{Status: model.SyncStatusStarted})

	allSucceeded := true

	for _, manager := range a.managers {
		events := make(chan model.SyncEvent)

		var syncError error

		go func() {
			defer func() {
				if panicValue := recover(); panicValue != nil {
					if err, ok := panicValue.(error); ok {
						syncError = err
					} else {
						syncError = errors.Errorf("sync failed: %s", panicValue)
					}
				}
				close(events)
			}()
			manager.Sync(events)
		}()

		for v := range events {
			if v.Status == model.SyncStatusAborted {
				allSucceeded = false
			}

			syncScreen.Send(sync.UpdateStateMsg{
				ManagerState: &sync.ManagerState{
					Key:    manager.Key(),
					Status: v.Status,
					Lines:  v.Lines,
					Input:  v.Login,
				},
			})
		}

		if syncError != nil {
			allSucceeded = false
			log.Error().Err(syncError).Msg("Uncaught panic while syncing")
		}
	}

	if allSucceeded {
		syncScreen.Send(sync.UpdateStateMsg{Status: model.SyncStatusFinished})
	} else {
		syncScreen.Send(sync.UpdateStateMsg{Status: model.SyncStatusAborted})
	}
}
