package app

import (
	"strings"

	"emperror.dev/errors"

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

	var err error
	doneChannel := make(chan struct{})

	go func() {
		defer close(doneChannel)
		err = a.startSyncManagers(syncScreen)
	}()

	syncScreen.Start()

	<-doneChannel
	if err != nil {
		panic(err)
	}
}

func (a *appImpl) startSyncManagers(syncScreen *sync.Screen) error {
	syncScreen.Send(sync.UpdateStateMsg{Status: model.SyncStatusStarted})

	for _, manager := range a.managers {
		events := make(chan model.SyncEvent)

		var syncError error
		syncDone := make(chan struct{})

		go func() {
			defer func() {
				if panicValue := recover(); panicValue != nil {
					if err, ok := panicValue.(error); ok {
						syncError = err
					} else {
						syncError = errors.Errorf("sync failed: %s", panicValue)
					}
				}
				close(syncDone)
				close(events)
			}()
			manager.Sync(events)
		}()

		for v := range events {
			syncScreen.Send(sync.UpdateStateMsg{
				ManagerState: &sync.ManagerState{
					Key:    manager.Key(),
					Status: v.Status,
					Lines:  v.Lines,
					Input:  v.Login,
				},
			})
		}

		<-syncDone
		if syncError != nil {
			return syncError
		}
	}

	syncScreen.Send(sync.UpdateStateMsg{Status: model.SyncStatusFinished})
	return nil
}
