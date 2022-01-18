package app

import (
	"io/ioutil"
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
	assert.NoError(t, ioutil.WriteFile(cfgFile, []byte("invalid"), 0o600))

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	_ = assertutil.AssertPanicsWithError(t, config.ErrInvalidConfig, func() {
		_ = NewApp(WithConfigService(config.NewService(config.WithViper(v))))
	})
}

func Test_NewAppNoManagers(t *testing.T) {
	cfg := config.VersionWrapper{Version: "1.0.0", Config: config.Config{}}
	cfgFile := path.Join(t.TempDir(), "temp-config.yaml")

	if cfgBytes, err := yaml.Marshal(cfg); err != nil {
		assert.NoError(t, err)
	} else {
		assert.NoError(t, ioutil.WriteFile(cfgFile, cfgBytes, 0o600))
	}

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	term := uiMocks.Terminal{}
	term.On("ApplyConfig", mock.AnythingOfType("ui.Config"), mock.Anything).Return()

	provider := managerMocks.Provider{}
	provider.On("CreateManager", mock.Anything, mock.Anything).Return([]managers.Manager{}, nil)

	app := NewApp(
		WithConfigService(config.NewService(config.WithViper(v))),
		WithTerminal(&term),
		WithProvider(&provider),
	)

	term.AssertNumberOfCalls(t, "ApplyConfig", 1)
	provider.AssertNumberOfCalls(t, "CreateManager", 1)

	assert.Len(t, app.(*appImpl).managers, 0)
}

func Test_appImpl_GetAllSnippets(t *testing.T) {
	snippets := []model.Snippet{
		{UUID: "uuid1", TitleFunc: testutil.FixedString("title-1"), LanguageFunc: testutil.FixedLanguage(model.LanguageYAML), TagUUIDs: []string{}, ContentFunc: testutil.FixedString("content-1")},
		{UUID: "uuid2", TitleFunc: testutil.FixedString("title-1"), LanguageFunc: testutil.FixedLanguage(model.LanguageBash), TagUUIDs: []string{}, ContentFunc: testutil.FixedString("content-2")},
	}

	manager := managerMocks.Manager{}
	manager.On("GetSnippets").Return(snippets, nil)

	app := appImpl{managers: []managers.Manager{&manager}}

	s := app.getAllSnippets()
	assertutil.AssertSnippetsEqual(t, snippets, s)
}
