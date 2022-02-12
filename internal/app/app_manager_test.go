package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/sync"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	managerMocks "github.com/lemoony/snipkit/mocks/managers"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
	syncMocks "github.com/lemoony/snipkit/mocks/ui/sync"
)

func Test_SyncManager(t *testing.T) {
	syncScreen := syncMocks.SyncScreen{}
	syncScreen.On("Start").Return()
	syncScreen.On("Send", mock.Anything).Return()

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowSync").Return(&syncScreen)

	managerSyncCloseChannel := make(chan time.Time)
	manager := managerMocks.Manager{}
	manager.On("Sync", mock.Anything).Return().WaitFor = managerSyncCloseChannel
	manager.On("Key").Return(model.ManagerKey("Manager X"))
	app := NewApp(WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManager(&manager))

	go func() {
		time.Sleep(100 * time.Millisecond)
		syncScreen.AssertCalled(t, "Start")
		syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{Status: model.SyncStatusStarted})

		time.Sleep(100 * time.Millisecond)
		manager.AssertCalled(t, "Sync", mock.Anything)
		call := manager.Calls[0]
		eventsChannel := call.Arguments.Get(0).(model.SyncEventChannel)
		assert.NotNil(t, eventsChannel)

		eventsChannel <- model.SyncEvent{
			Status: model.SyncStatusStarted,
			Lines:  []model.SyncLine{{Type: model.SyncLineTypeInfo, Value: "Manager X started"}},
		}

		time.Sleep(time.Millisecond * 100)

		eventsChannel <- model.SyncEvent{
			Status: model.SyncStatusFinished,
			Lines:  []model.SyncLine{{Type: model.SyncLineTypeInfo, Value: "Manager X finished"}},
		}

		close(managerSyncCloseChannel)
	}()

	app.SyncManager()

	syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{Status: model.SyncStatusStarted})
	syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{
		ManagerState: &sync.ManagerState{
			Key:    manager.Key(),
			Status: model.SyncStatusStarted,
			Lines:  []model.SyncLine{{Type: model.SyncLineTypeInfo, Value: "Manager X started"}},
		},
	})
	syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{
		ManagerState: &sync.ManagerState{
			Key:    manager.Key(),
			Status: model.SyncStatusFinished,
			Lines:  []model.SyncLine{{Type: model.SyncLineTypeInfo, Value: "Manager X finished"}},
		},
	})

	syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{Status: model.SyncStatusFinished})
}

func Test_Sync_manager_panic(t *testing.T) {
	syncScreen := syncMocks.SyncScreen{}
	syncScreen.On("Start").Return()
	syncScreen.On("Send", mock.Anything).Return()

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowSync").Return(&syncScreen)

	managerSyncCloseChannel := make(chan time.Time)
	manager := managerMocks.Manager{}
	manager.On("Sync", mock.Anything).Panic("test panic").WaitFor = managerSyncCloseChannel
	manager.On("Key").Return(model.ManagerKey("Manager X"))
	app := NewApp(WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManager(&manager))

	go func() {
		time.Sleep(100 * time.Millisecond)
		syncScreen.AssertCalled(t, "Start")
		syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{Status: model.SyncStatusStarted})
		close(managerSyncCloseChannel)
	}()

	assert.Panics(t, func() {
		app.SyncManager()
	})

	syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{Status: model.SyncStatusStarted})
}
