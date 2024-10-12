package app

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/assertutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	managerMocks "github.com/lemoony/snipkit/mocks/managers"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_NewApp_NoConfigFile(t *testing.T) {
	v := viper.NewWithOptions()

	_ = assertutil.AssertPanicsWithError(t, config.ErrConfigNotFound{}, func() {
		_ = NewApp(WithConfigService(config.NewService(config.WithViper(v))))
	})
}

func Test_NewAppInvalidConfigFile(t *testing.T) {
	cfgFile := path.Join(t.TempDir(), "invalid-config")
	assert.NoError(t, os.WriteFile(cfgFile, []byte("invalid"), 0o600))

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	_ = assertutil.AssertPanicsWithError(t, config.ErrInvalidConfig, func() {
		_ = NewApp(WithConfigService(config.NewService(config.WithViper(v))))
	})
}

func Test_NewAppNeedsConfigMigration(t *testing.T) {
	cfg := config.VersionWrapper{Version: "1.1.0", Config: config.Config{}}
	cfgFile := path.Join(t.TempDir(), "temp-config.yaml")

	if cfgBytes, err := yaml.Marshal(cfg); err != nil {
		assert.NoError(t, err)
	} else {
		assert.NoError(t, os.WriteFile(cfgFile, cfgBytes, 0o600))
	}

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.AnythingOfType("ui.Config"), mock.Anything).Return()

	err := assertutil.AssertPanicsWithError(t, ErrMigrateConfig{}, func() {
		_ = NewApp(
			WithConfigService(config.NewService(config.WithViper(v))),
			WithTUI(&tui),
		)
	}).(ErrMigrateConfig)

	assert.Equal(t, config.Version, err.latestVersion)
	assert.Equal(t, "1.1.0", err.currentVersion)
	assert.Contains(t, err.Error(), "to migrate the config file from version 1.1.0 to 1.1.2.")
}

func Test_NewAppNoManagers(t *testing.T) {
	cfg := config.VersionWrapper{Version: config.Version, Config: config.Config{}}
	cfgFile := path.Join(t.TempDir(), "temp-config.yaml")

	if cfgBytes, err := yaml.Marshal(cfg); err != nil {
		assert.NoError(t, err)
	} else {
		assert.NoError(t, os.WriteFile(cfgFile, cfgBytes, 0o600))
	}

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.AnythingOfType("ui.Config"), mock.Anything).Return()

	provider := managerMocks.Provider{}
	provider.On("CreateManager", mock.Anything, mock.Anything).Return([]managers.Manager{}, nil)

	app := NewApp(
		WithConfigService(config.NewService(config.WithViper(v))),
		WithTUI(&tui),
		WithProvider(&provider),
	)

	tui.AssertNumberOfCalls(t, mockutil.ApplyConfig, 1)
	provider.AssertNumberOfCalls(t, "CreateManager", 1)

	assert.Len(t, app.(*appImpl).managers, 0)
}

func Test_appImpl_GetAllSnippets(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid1", Title: "title-1", Language: model.LanguageYAML, Tags: []string{}, Content: "content-1"},
		testutil.TestSnippet{ID: "uuid2", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: "content-2"},
	}

	manager := managerMocks.Manager{}
	manager.On("GetSnippets").Return(snippets, nil)

	app := appImpl{managers: []managers.Manager{&manager}}

	s := app.getAllSnippets()
	assertutil.AssertSnippetsEqual(t, snippets, s)
}
