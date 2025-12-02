package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/githubgist"
	"github.com/lemoony/snipkit/internal/managers/masscode"
	"github.com/lemoony/snipkit/internal/managers/pet"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/sync"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	configMocks "github.com/lemoony/snipkit/mocks/config"
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

	app.SyncManager()

	syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{Status: model.SyncStatusStarted})
	syncScreen.AssertCalled(t, "Send", sync.UpdateStateMsg{Status: model.SyncStatusAborted})
}

func Test_AddManager_UserConfirmsAddition(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	// User selects first manager
	tui.On(mockutil.ShowPicker, "Which snippet manager should be added to your configuration", mock.Anything, mock.Anything).Return(0, true)

	// User confirms the change
	tui.On(mockutil.Confirmation, mock.Anything).Return(true)
	tui.On(mockutil.Print, mock.Anything).Return()

	cfg := configtest.NewTestConfig().Config

	provider := managerMocks.Provider{}
	provider.On("ManagerDescriptions", cfg.Manager).Return([]model.ManagerDescription{
		{Key: fslibrary.Key, Name: "FS Library", Description: "File system library"},
	})
	// AutoConfig returns a config with FsLibrary set
	provider.On("AutoConfig", fslibrary.Key, mock.Anything).Return(managers.Config{
		FsLibrary: &fslibrary.Config{Enabled: true, LibraryPath: []string{"/test/path"}},
	})
	provider.On("CreateManager", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]managers.Manager{}, nil)

	configService := configMocks.ConfigService{}
	configService.On("UpdateManagerConfig", mock.Anything).Return()
	configService.On("ConfigFilePath").Return("/test/config.yml")
	configService.On("LoadConfig").Return(cfg, nil)
	configService.On("NeedsMigration").Return(false, "")

	app := NewApp(
		WithTUI(&tui),
		WithProvider(&provider),
		WithConfigService(&configService),
	)

	app.AddManager()

	// Verify the flow
	tui.AssertCalled(t, mockutil.ShowPicker, "Which snippet manager should be added to your configuration", mock.Anything, mock.Anything)

	// Verify Confirmation was called
	tui.AssertCalled(t, mockutil.Confirmation, mock.Anything)

	// Verify config was updated
	configService.AssertCalled(t, "UpdateManagerConfig", mock.MatchedBy(func(cfg managers.Config) bool {
		return cfg.FsLibrary != nil && cfg.FsLibrary.Enabled
	}))

	// Verify success message was printed
	tui.AssertCalled(t, mockutil.Print, mock.Anything)
}

func Test_AddManager_UserDeclinesAddition(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowPicker, mock.Anything, mock.Anything, mock.Anything).Return(0, true)
	tui.On(mockutil.Confirmation, mock.Anything).Return(false) // User declines
	tui.On(mockutil.Print, mock.Anything).Return()

	cfg := configtest.NewTestConfig().Config

	provider := managerMocks.Provider{}
	provider.On("ManagerDescriptions", cfg.Manager).Return([]model.ManagerDescription{
		{Key: pet.Key, Name: "Pet", Description: "Pet manager"},
	})
	provider.On("AutoConfig", pet.Key, mock.Anything).Return(managers.Config{
		Pet: &pet.Config{Enabled: true},
	})
	provider.On("CreateManager", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]managers.Manager{}, nil)

	configService := configMocks.ConfigService{}
	configService.On("ConfigFilePath").Return("/test/config.yml")
	configService.On("LoadConfig").Return(cfg, nil)
	configService.On("NeedsMigration").Return(false, "")

	app := NewApp(
		WithTUI(&tui),
		WithProvider(&provider),
		WithConfigService(&configService),
	)

	app.AddManager()

	// Verify UpdateManagerConfig was NOT called
	configService.AssertNotCalled(t, "UpdateManagerConfig")

	// But Print was called
	tui.AssertCalled(t, mockutil.Print, mock.Anything)
}

func Test_AddManager_UserCancelsSelection(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowPicker, mock.Anything, mock.Anything, mock.Anything).Return(0, false) // User cancels

	cfg := configtest.NewTestConfig().Config

	provider := managerMocks.Provider{}
	provider.On("ManagerDescriptions", cfg.Manager).Return([]model.ManagerDescription{
		{Key: fslibrary.Key, Name: "FS Library", Description: "File system library"},
	})
	provider.On("CreateManager", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]managers.Manager{}, nil)

	app := NewApp(
		WithTUI(&tui),
		WithConfig(cfg),
		WithProvider(&provider),
	)

	app.AddManager()

	// Verify picker was shown
	tui.AssertCalled(t, mockutil.ShowPicker, mock.Anything, mock.Anything, mock.Anything)

	// Verify no further calls were made
	tui.AssertNotCalled(t, mockutil.Confirmation, mock.Anything)
	tui.AssertNotCalled(t, mockutil.Print, mock.Anything)
}

func Test_AddManager_AllManagerTypes(t *testing.T) {
	tests := getAddManagerTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddManagerType(t, tt.managerKey, tt.managerName, tt.configProvider())
		})
	}
}

func getAddManagerTestCases() []struct {
	name           string
	managerKey     model.ManagerKey
	managerName    string
	configProvider func() managers.Config
} {
	return []struct {
		name           string
		managerKey     model.ManagerKey
		managerName    string
		configProvider func() managers.Config
	}{
		{"FsLibrary", fslibrary.Key, "FS Library", func() managers.Config {
			return managers.Config{FsLibrary: &fslibrary.Config{Enabled: true, LibraryPath: []string{"/test"}}}
		}},
		{"SnippetsLab", snippetslab.Key, "SnippetsLab", func() managers.Config {
			return managers.Config{SnippetsLab: &snippetslab.Config{Enabled: true}}
		}},
		{"PictarineSnip", pictarinesnip.Key, "Pictarine Snip", func() managers.Config {
			return managers.Config{PictarineSnip: &pictarinesnip.Config{Enabled: true}}
		}},
		{"Pet", pet.Key, "Pet", func() managers.Config {
			return managers.Config{Pet: &pet.Config{Enabled: true}}
		}},
		{"MassCode", masscode.Key, "massCode", func() managers.Config {
			return managers.Config{MassCode: &masscode.Config{Enabled: true}}
		}},
		{"GithubGist", githubgist.Key, "GitHub Gist", func() managers.Config {
			return managers.Config{GithubGist: &githubgist.Config{Enabled: true}}
		}},
	}
}

func testAddManagerType(t *testing.T, managerKey model.ManagerKey, managerName string, managerConfig managers.Config) {
	t.Helper()
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowPicker, mock.Anything, mock.Anything, mock.Anything).Return(0, true)
	tui.On(mockutil.Confirmation, mock.Anything).Return(true)
	tui.On(mockutil.Print, mock.Anything).Return()

	cfg := configtest.NewTestConfig().Config

	provider := managerMocks.Provider{}
	provider.On("ManagerDescriptions", cfg.Manager).Return([]model.ManagerDescription{
		{Key: managerKey, Name: managerName, Description: "test"},
	})
	provider.On("AutoConfig", managerKey, mock.Anything).Return(managerConfig)
	provider.On("CreateManager", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]managers.Manager{}, nil)

	configService := configMocks.ConfigService{}
	configService.On("UpdateManagerConfig", mock.Anything).Return()
	configService.On("ConfigFilePath").Return("/test/config.yml")
	configService.On("LoadConfig").Return(cfg, nil)
	configService.On("NeedsMigration").Return(false, "")

	app := NewApp(
		WithTUI(&tui),
		WithProvider(&provider),
		WithConfigService(&configService),
	)

	app.AddManager()

	// Verify the manager config was updated
	configService.AssertCalled(t, "UpdateManagerConfig", mock.Anything)

	// Verify confirmation was shown (which means the diff was computed)
	tui.AssertCalled(t, mockutil.Confirmation, mock.Anything)
	if call := mockutil.FindMethodCall(mockutil.Confirmation, tui.Calls); call != nil {
		// Verify the confirmation message contains the diff rendering
		assert.NotNil(t, call.Arguments.Get(0))
	}
}
